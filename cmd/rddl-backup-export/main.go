package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	claimPrefix          = "Claim/"
	confirmedClaimPrefix = "ConfirmedClaim/"
	countKey             = "Count"
)

type RedeemClaim struct {
	ID           int    `json:"id"`
	Beneficiary  string `json:"beneficiary"`
	Amount       uint64 `json:"amount"`
	LiquidTXHash string `json:"liquid-tx-hash"`
	ClaimID      int    `json:"claim-id"`
}

type row struct {
	KeyID int
	Claim RedeemClaim
}

func decodeID(key []byte, prefix string) (int, error) {
	if !bytes.HasPrefix(key, []byte(prefix)) {
		return 0, fmt.Errorf("missing prefix %q", prefix)
	}

	idBytes := key[len(prefix):]
	if len(idBytes) != 8 {
		return 0, fmt.Errorf("unexpected id length %d", len(idBytes))
	}

	return int(binary.BigEndian.Uint64(idBytes)), nil
}

func mdEscape(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func main() {
	dbPath := flag.String("db", "data", "path to LevelDB directory")
	outPath := flag.String("out", "data/backup_export.md", "path to output Markdown file")
	flag.Parse()

	db, err := leveldb.OpenFile(*dbPath, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open db: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	var count string
	var unconfirmed []row
	var confirmed []row
	var unknown []string

	for iter.Next() {
		k := append([]byte(nil), iter.Key()...)
		v := append([]byte(nil), iter.Value()...)
		ks := string(k)

		switch {
		case ks == countKey:
			count = string(v)

		case strings.HasPrefix(ks, claimPrefix):
			id, derr := decodeID(k, claimPrefix)
			if derr != nil {
				unknown = append(unknown, fmt.Sprintf("Claim key decode error: %v", derr))
				continue
			}

			var rc RedeemClaim
			if err := json.Unmarshal(v, &rc); err != nil {
				unknown = append(unknown, fmt.Sprintf("Claim/%d json error: %v", id, err))
				continue
			}
			unconfirmed = append(unconfirmed, row{KeyID: id, Claim: rc})

		case strings.HasPrefix(ks, confirmedClaimPrefix):
			id, derr := decodeID(k, confirmedClaimPrefix)
			if derr != nil {
				unknown = append(unknown, fmt.Sprintf("ConfirmedClaim key decode error: %v", derr))
				continue
			}

			var rc RedeemClaim
			if err := json.Unmarshal(v, &rc); err != nil {
				unknown = append(unknown, fmt.Sprintf("ConfirmedClaim/%d json error: %v", id, err))
				continue
			}
			confirmed = append(confirmed, row{KeyID: id, Claim: rc})

		default:
			unknown = append(unknown, fmt.Sprintf("Unknown key: %q (%d bytes)", ks, len(v)))
		}
	}

	if err := iter.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "iterate db: %v\n", err)
		os.Exit(1)
	}

	sort.Slice(unconfirmed, func(i, j int) bool { return unconfirmed[i].KeyID < unconfirmed[j].KeyID })
	sort.Slice(confirmed, func(i, j int) bool { return confirmed[i].KeyID < confirmed[j].KeyID })

	f, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create output: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintln(f, "# RDDL Claim Service Backup Export")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "- Database path: `%s`\n", mdEscape(*dbPath))
	if count != "" {
		fmt.Fprintf(f, "- Count key: `%s`\n", mdEscape(count))
	} else {
		fmt.Fprintln(f, "- Count key: _missing_")
	}
	fmt.Fprintf(f, "- Unconfirmed claims: `%d`\n", len(unconfirmed))
	fmt.Fprintf(f, "- Confirmed claims: `%d`\n", len(confirmed))
	fmt.Fprintln(f)

	fmt.Fprintln(f, "## Unconfirmed Claims")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "| KeyID | ClaimID | Beneficiary | Amount | LiquidTXHash | JSON.ID |")
	fmt.Fprintln(f, "|---:|---:|---|---:|---|---:|")
	for _, r := range unconfirmed {
		fmt.Fprintf(f, "| %d | %d | %s | %d | %s | %d |\n",
			r.KeyID,
			r.Claim.ClaimID,
			mdEscape(r.Claim.Beneficiary),
			r.Claim.Amount,
			mdEscape(r.Claim.LiquidTXHash),
			r.Claim.ID,
		)
	}
	if len(unconfirmed) == 0 {
		fmt.Fprintln(f, "| _none_ |  |  |  |  |  |")
	}
	fmt.Fprintln(f)

	fmt.Fprintln(f, "## Confirmed Claims")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "| KeyID | ClaimID | Beneficiary | Amount | LiquidTXHash | JSON.ID |")
	fmt.Fprintln(f, "|---:|---:|---|---:|---|---:|")
	for _, r := range confirmed {
		fmt.Fprintf(f, "| %d | %d | %s | %d | %s | %d |\n",
			r.KeyID,
			r.Claim.ClaimID,
			mdEscape(r.Claim.Beneficiary),
			r.Claim.Amount,
			mdEscape(r.Claim.LiquidTXHash),
			r.Claim.ID,
		)
	}
	if len(confirmed) == 0 {
		fmt.Fprintln(f, "| _none_ |  |  |  |  |  |")
	}

	if len(unknown) > 0 {
		fmt.Fprintln(f)
		fmt.Fprintln(f, "## Unknown Or Decode Errors")
		fmt.Fprintln(f)
		for _, u := range unknown {
			fmt.Fprintf(f, "- %s\n", mdEscape(u))
		}
	}

	fmt.Printf("wrote %s\n", *outPath)
}
