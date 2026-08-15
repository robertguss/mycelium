package check

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/robertguss/mycelium/internal/ladder"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/pack"
)

var lowercaseHex64 = regexp.MustCompile(`^[0-9a-f]{64}$`)

const (
	contractCMP = "program/packs/council/contracts/commissioning.md"
	contractRPT = "program/packs/council/contracts/report.md"
	contractRCL = "program/packs/council/contracts/reconciliation.md"
)

type ladderCMP struct {
	id      string
	rung    ladder.Rung
	adapter ladder.Adapter
	hash    string
}

type ladderRPT struct {
	id            string
	commissioning string
	rung          ladder.Rung
	adapter       ladder.Adapter
	promptSHA     string
	dissent       string
}

type ladderRCL struct {
	id            string
	commissioning string
	rung          ladder.Rung
	retained      string
}

// checkLadder binds §6 IFF rules when the council pack is present (items 20–22).
func checkLadder(root string, arts []artifactFile, packs []pack.Pack, add func(string, string, string, string)) {
	if !pack.ReviewsAllowed(packs) {
		return
	}

	cmps := map[string]ladderCMP{}
	var rpts []ladderRPT
	var rcls []ladderRCL

	for _, a := range arts {
		switch a.Home {
		case "reviews/commissioning":
			c, ok := loadLadderCMP(root, a, add)
			if ok {
				cmps[c.id] = c
			}
		case "reviews/reports":
			r, ok := loadLadderRPT(root, a)
			if ok {
				rpts = append(rpts, r)
			}
		case "reviews/reconciliations":
			r, ok := loadLadderRCL(root, a, add)
			if ok {
				rcls = append(rcls, r)
			}
		}
	}

	rptsByCMP := map[string][]ladderRPT{}
	for _, r := range rpts {
		if _, ok := cmps[r.commissioning]; !ok {
			add(
				fmt.Sprintf("%s commissioning %q does not resolve", r.id, r.commissioning),
				"commissioning-resolve",
				contractRPT,
				"set commissioning to an existing CMP-###, or create that commissioning file",
			)
			continue
		}
		rptsByCMP[r.commissioning] = append(rptsByCMP[r.commissioning], r)
	}

	rclsByCMP := map[string][]ladderRCL{}
	for _, r := range rcls {
		if _, ok := cmps[r.commissioning]; !ok {
			add(
				fmt.Sprintf("%s commissioning %q does not resolve", r.id, r.commissioning),
				"commissioning-resolve",
				contractRCL,
				"set commissioning to an existing CMP-###, or create that commissioning file",
			)
			continue
		}
		rclsByCMP[r.commissioning] = append(rclsByCMP[r.commissioning], r)
	}

	for id, c := range cmps {
		matchingRPT := rptsByCMP[id]
		matchingRCL := rclsByCMP[id]

		if err := ladder.Cardinality(c.rung, len(matchingRPT), len(matchingRCL)); err != nil {
			switch c.rung {
			case ladder.RungCouncil:
				add(
					fmt.Sprintf("%s council requires >=2 model reports and exactly one reconciliation", id),
					"council-cardinality",
					contractCMP,
					"add matching RPT-### files (>=2) and exactly one RCL-###, or remove started reports",
				)
			case ladder.RungSecondOpinion:
				add(
					fmt.Sprintf("%s second-opinion requires exactly one model report and no reconciliation", id),
					"second-opinion-cardinality",
					contractCMP,
					"keep exactly one matching RPT-### and no RCL-###, or remove started reports",
				)
			}
		}

		for _, r := range matchingRPT {
			if r.rung != c.rung {
				add(
					fmt.Sprintf("%s rung %q does not match commissioning %s rung %q", r.id, r.rung, id, c.rung),
					"prompt-identity",
					contractRPT,
					fmt.Sprintf("set rung to %q to match %s", c.rung, id),
				)
			}
			if r.adapter != c.adapter {
				add(
					fmt.Sprintf("%s adapter %q does not match commissioning %s adapter %q", r.id, r.adapter, id, c.adapter),
					"prompt-identity",
					contractRPT,
					fmt.Sprintf("set adapter to %q to match %s", c.adapter, id),
				)
			}
			if !lowercaseHex64.MatchString(r.promptSHA) || r.promptSHA != c.hash {
				add(
					fmt.Sprintf("%s prompt_sha256 mismatch (want %s)", r.id, c.hash),
					"prompt-identity",
					contractRPT,
					fmt.Sprintf("set prompt_sha256 to the sha256 hex of %s ## Prompt (trim surrounding whitespace)", id),
				)
			}
		}

		if c.rung == ladder.RungCouncil && len(matchingRCL) == 1 {
			var dissentBodies []string
			for _, r := range matchingRPT {
				dissentBodies = append(dissentBodies, r.dissent)
			}
			rcl := matchingRCL[0]
			if err := ladder.SeedDissentOK(dissentBodies, rcl.retained); err != nil {
				add(
					fmt.Sprintf("%s ## Retained dissent missing SEED-DISSENT", rcl.id),
					"seeded-dissent",
					contractRCL,
					"retain the SEED-DISSENT token in ## Retained dissent",
				)
			}
		}
	}
}

func loadLadderCMP(root string, a artifactFile, add func(string, string, string, string)) (ladderCMP, bool) {
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.Rel)))
	if err != nil {
		return ladderCMP{}, false
	}
	doc, err := metadata.Parse(b)
	if err != nil {
		return ladderCMP{}, false
	}

	if !ladder.OptInOK(doc.Meta["opt_in"]) {
		got := formatOptIn(doc.Meta["opt_in"])
		add(
			fmt.Sprintf("%s opt_in must be true (got %s)", a.IDStr, got),
			"council-opt-in",
			contractCMP,
			"set opt_in = true, or delete the commissioning file",
		)
	}

	rungStr, _ := doc.Meta["rung"].(string)
	rung, err := ladder.ParseRung(rungStr)
	if err != nil {
		return ladderCMP{}, false
	}
	classStr, _ := doc.Meta["cost_class"].(string)
	class, err := ladder.ParseCostClass(classStr)
	if err != nil {
		return ladderCMP{}, false
	}
	if !ladder.CostClassOK(rung, class) {
		if rung == ladder.RungCouncil {
			add(
				fmt.Sprintf("%s cost_class %q is not quick|standard|high-stakes when rung=council", a.IDStr, classStr),
				"council-cost-class",
				contractCMP,
				"set cost_class to quick, standard, or high-stakes",
			)
		} else {
			add(
				fmt.Sprintf("%s cost_class %q is not cheap when rung=second-opinion", a.IDStr, classStr),
				"council-cost-class",
				contractCMP,
				"set cost_class to cheap",
			)
		}
	}

	adapterStr, _ := doc.Meta["adapter"].(string)
	adapter, err := ladder.ParseAdapter(adapterStr)
	if err != nil {
		return ladderCMP{}, false
	}

	return ladderCMP{
		id:      a.IDStr,
		rung:    rung,
		adapter: adapter,
		hash:    ladder.PromptSHA256(doc.Body),
	}, true
}

func loadLadderRPT(root string, a artifactFile) (ladderRPT, bool) {
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.Rel)))
	if err != nil {
		return ladderRPT{}, false
	}
	doc, err := metadata.Parse(b)
	if err != nil {
		return ladderRPT{}, false
	}
	comm, _ := doc.Meta["commissioning"].(string)
	rungStr, _ := doc.Meta["rung"].(string)
	rung, err := ladder.ParseRung(rungStr)
	if err != nil {
		return ladderRPT{}, false
	}
	adapterStr, _ := doc.Meta["adapter"].(string)
	adapter, err := ladder.ParseAdapter(adapterStr)
	if err != nil {
		return ladderRPT{}, false
	}
	sha, _ := doc.Meta["prompt_sha256"].(string)
	return ladderRPT{
		id:            a.IDStr,
		commissioning: comm,
		rung:          rung,
		adapter:       adapter,
		promptSHA:     sha,
		dissent:       ladder.SectionBody(doc.Body, "Dissent"),
	}, true
}

func loadLadderRCL(root string, a artifactFile, add func(string, string, string, string)) (ladderRCL, bool) {
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.Rel)))
	if err != nil {
		return ladderRCL{}, false
	}
	doc, err := metadata.Parse(b)
	if err != nil {
		return ladderRCL{}, false
	}
	comm, _ := doc.Meta["commissioning"].(string)
	rungStr, _ := doc.Meta["rung"].(string)
	rung, err := ladder.ParseRung(rungStr)
	if err != nil {
		add(
			fmt.Sprintf("%s rung must be council (got %q)", a.IDStr, rungStr),
			"council-cardinality",
			contractRCL,
			"set rung = \"council\", or delete the reconciliation file",
		)
		return ladderRCL{}, false
	}
	if rung != ladder.RungCouncil {
		add(
			fmt.Sprintf("%s rung must be council (got %q)", a.IDStr, rung),
			"council-cardinality",
			contractRCL,
			"set rung = \"council\", or delete the reconciliation file",
		)
	}
	return ladderRCL{
		id:            a.IDStr,
		commissioning: comm,
		rung:          rung,
		retained:      ladder.SectionBody(doc.Body, "Retained dissent"),
	}, true
}

func formatOptIn(v any) string {
	switch x := v.(type) {
	case nil:
		return "missing"
	case bool:
		return fmt.Sprintf("%v", x)
	case string:
		return fmt.Sprintf("%q", x)
	default:
		return fmt.Sprintf("%v", x)
	}
}
