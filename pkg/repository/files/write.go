package repository

import (
	"os"
	"path"

	"google.golang.org/protobuf/encoding/protojson"
)

func (r *repository) writeFile(pth string, buf []byte) error {
	pth = path.Join(r.url, pth)
	err := os.MkdirAll(path.Dir(pth), os.ModeDir)
	if err != nil {
		return err
	}
	f, err := os.Create(pth)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) writeConfig() error {
	buf, err := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
	}.Marshal(r.config)
	if err != nil {
		return err
	}

	err = r.writeFile(pathConfig, buf)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) writeAccounts() error {
	buf, err := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
	}.Marshal(r.accounts)
	if err != nil {
		return err
	}

	err = r.writeFile(pathAccounts, buf)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) writeClients() error {
	buf, err := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
	}.Marshal(r.clients)
	if err != nil {
		return err
	}

	err = r.writeFile(pathClients, buf)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) writeTemplates() error {
	for _, tpl := range r.templates {
		buf, err := protojson.MarshalOptions{
			Indent:          "  ",
			EmitUnpopulated: true,
		}.Marshal(tpl)
		if err != nil {
			return err
		}

		err = r.writeFile(path.Join(pathTemplates, tpl.FileName), buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) writeProjects() error {
	for _, prj := range r.projects {
		buf, err := protojson.MarshalOptions{
			Indent:          "  ",
			EmitUnpopulated: true,
		}.Marshal(prj)
		if err != nil {
			return err
		}

		err = r.writeFile(path.Join(pathProjects, prj.FileName), buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) writeInvoices() error {
	for _, inv := range r.invoices {
		buf, err := protojson.MarshalOptions{
			Indent:          "  ",
			EmitUnpopulated: true,
		}.Marshal(inv)
		if err != nil {
			return err
		}

		err = r.writeFile(path.Join(pathInvoices, inv.Year, inv.FileName), buf)
		if err != nil {
			return err
		}
	}
	return nil
}
