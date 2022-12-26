package repository

import (
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/itohio/phingo/pkg/types"
	"google.golang.org/protobuf/encoding/protojson"
)

func (r *repository) readFile(pth string) ([]byte, error) {
	f, err := r.fs.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (r *repository) readConfig() error {
	buf, err := r.readFile(pathConfig)
	if err != nil {
		return err
	}
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(buf, r.config)
}

func (r *repository) readAccounts() error {
	buf, err := r.readFile(pathAccounts)
	if err != nil {
		return err
	}
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(buf, r.accounts)
}

func (r *repository) readClients() error {
	buf, err := r.readFile(pathClients)
	if err != nil {
		return err
	}
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(buf, r.clients)
}

func (r *repository) readTemplates() error {
	return fs.WalkDir(
		r.fs,
		pathTemplates,
		func(pth string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(pth, ".md") {
				return nil
			}
			buf, err := r.readFile(pth)
			if err != nil {
				return err
			}

			id := path.Base(pth)
			id = id[:len(id)-3]
			tpl := &types.Template{
				Id:       id,
				What:     id,
				FileName: pth,
				Text:     buf,
			}
			r.templates = append(r.templates, tpl)
			return nil
		},
	)
}

func (r *repository) readProjects() error {
	return fs.WalkDir(
		r.fs,
		pathProjects,
		func(pth string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(pth, ".json") {
				return nil
			}
			buf, err := r.readFile(pth)
			if err != nil {
				return err
			}
			var prj types.Project
			err = protojson.UnmarshalOptions{
				DiscardUnknown: true,
				AllowPartial:   true,
			}.Unmarshal(buf, &prj)
			if err != nil {
				return err
			}
			prj.FileName = pth
			id := path.Base(pth)
			id = id[:len(id)-5]
			if prj.Id == "" {
				prj.Id = id
			}
			if prj.Name == "" {
				prj.Name = strings.ToTitle(id)
			}
			r.projects = append(r.projects, &prj)
			return nil
		},
	)
}

func (r *repository) readInvoices() error {
	return fs.WalkDir(
		r.fs,
		pathInvoices,
		func(pth string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(pth, ".json") {
				return nil
			}
			b, id := path.Dir(pth), path.Base(pth)
			if path.Dir(b) != pathInvoices {
				return nil
			}

			buf, err := r.readFile(pth)
			if err != nil {
				return err
			}
			var inv types.Invoice
			err = protojson.UnmarshalOptions{
				DiscardUnknown: true,
				AllowPartial:   true,
			}.Unmarshal(buf, &inv)
			if err != nil {
				return err
			}
			inv.Year = path.Dir(b)
			inv.FileName = pth
			if inv.Id == "" {
				id = id[:len(id)-5]
				inv.Id = id
			}
			r.invoices = append(r.invoices, &inv)
			return nil
		},
	)
}
