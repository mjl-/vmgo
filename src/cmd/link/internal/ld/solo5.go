package ld

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"io"
	"strings"
)

// TODO(mjl): make less ugly. there's probably a better way to map c structures to bytes. perhaps just make them go structs and pack them to bytes?
// TODO(mjl): add the manifest.json to build dependencies. so we rebuild properly after changing it. for now i'm building with -a, will get annoying fast.

type writer struct {
	out io.Writer
	err error
}

func (w *writer) Put(v interface{}) {
	if w.err == nil {
		w.err = binary.Write(w.out, binary.LittleEndian, v)
	}
}

// Parse solo5 manifest.json, writing the contents of the ".note.solo5.manifest" elf section to ctxt.solo5Manifest.
func parseSolo5Manifest(ctx *Link) {
	f, err := os.Open(*solo5Manifest)
	if err != nil {
		Exitf("open solo5 manifest: %v", err)
	}
	defer f.Close()

	var manifest struct {
		Type    string `json:"type"`    // "solo5.manifest"
		Version int    `json:"version"` // 1
		Devices []struct {
			Name string `json:"name"`
			Type string `json:"type"` // "NET_BASIC" or "BLOCK_BASIC"
		} `json:"devices"`
	}
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	err = dec.Decode(&manifest)
	if err != nil {
		Exitf("decoding solo5 manifest JSON: %v", err)
	}
	if manifest.Type != "solo5.manifest" {
		Exitf("unknown solo5 manifest type %q, expected %q", manifest.Type, "solo5.manifest")
	}
	if manifest.Version != 1 {
		Exitf("unknown solo5 manifest version %d, expected %d", manifest.Version, 1)
	}

	mftbuf := &bytes.Buffer{}
	w := &writer{out: mftbuf}
	const mft_version = 1
	w.Put(uint32(mft_version))
	w.Put(uint32(1 + len(manifest.Devices)))

	// empty reserved first
	w.Put([]byte(strings.Repeat("\000", 68)))
	w.Put(uint32(1 << 30))
	w.Put([]byte(strings.Repeat("\000", (16+8+1+7) & ^7)))

	seen := map[string]struct{}{}
	for _, dev := range manifest.Devices {
		if dev.Name == "" {
			Exitf("solo5 manifest: empty device name")
		}
		if _, ok := seen[dev.Name]; ok {
			Exitf("solo5 manifest: duplicate device name %q", dev.Name)
		}
		seen[dev.Name] = struct{}{}

		w.Put([]byte(dev.Name))
		w.Put([]byte(strings.Repeat("\000", 68-len(dev.Name))))

		const mft_dev_block_basic = 1
		const mft_dev_net_basic = 2
		switch dev.Type {
		case "BLOCK_BASIC":
			w.Put(uint32(mft_dev_block_basic))
		case "NET_BASIC":
			w.Put(uint32(mft_dev_net_basic))
		default:
			Exitf("solo5 manifest: unknown device type %q", dev.Type)
		}
		w.Put([]byte(strings.Repeat("\000", (16+8+1+7) & ^7)))
	}

	if w.err != nil {
		Exitf("making solo5 manifest section: %v", w.err)
	}

	ctx.solo5Manifest = mftbuf.Bytes()
}
