package repository

import (
	"io/fs"
	"log"
	"path"

	"github.com/golang/protobuf/proto"
	defaultRepo "github.com/itohio/phingo/pkg/repository/default"
	"github.com/itohio/phingo/pkg/types"
	"google.golang.org/protobuf/encoding/protojson"
)

func (r *repository) writeFile(pth string, buf []byte) error {
	// FIXME: Make it more secure
	pth = path.Join(r.url, pth)
	err := r.fs.MkDirAll(path.Dir(pth), fs.ModeDir)
	if err != nil {
		return err
	}
	f, err := r.fs.Create(pth)
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

func (r *repository) removeFile(pth string) error {
	// FIXME: make it more secure
	pth = path.Join(r.url, pth)
	return r.fs.Rename(pth, pth+".rm")
}

func (r *repository) writeConfig() error {
	if !r.configModified {
		return nil
	}
	buf, err := protojson.Marshal(r.config)
	if err != nil {
		return err
	}

	buf, err = json2yaml(buf)
	if err != nil {
		return err
	}

	err = r.writeFile(defaultRepo.PathConfigYaml, buf)
	if err != nil {
		return err
	}
	r.configModified = false
	return nil
}

func (r *repository) writeAccounts() error {
	if !r.accountsModified {
		return nil
	}
	buf, err := protojson.Marshal(r.accounts)
	if err != nil {
		return err
	}

	buf, err = json2yaml(buf)
	if err != nil {
		return err
	}

	err = r.writeFile(defaultRepo.PathAccountsYaml, buf)
	if err != nil {
		return err
	}
	r.accountsModified = false
	return nil
}

func (r *repository) writeClients() error {
	if !r.clientsModified {
		return nil
	}
	buf, err := protojson.Marshal(r.clients)
	if err != nil {
		return err
	}

	buf, err = json2yaml(buf)
	if err != nil {
		return err
	}

	err = r.writeFile(defaultRepo.PathClientsYaml, buf)
	if err != nil {
		return err
	}
	r.clientsModified = false
	return nil
}

func (r *repository) writeTemplates() error {
	for _, tpl := range r.templates {
		if tpl.Id == "" {
			log.Fatalln("template has empty Id")
			continue
		}

		_, ok := r.templatesModified[tpl.Id]
		if !ok {
			continue
		}

		if tpl.FileName == "" {
			tpl.FileName = tpl.Id + ".md"
		}
		err := r.writeFile(path.Join(defaultRepo.PathTemplates, path.Base(tpl.FileName)), tpl.Text)
		if err != nil {
			return err
		}
	}

	for id, ms := range r.templatesModified {
		delete(r.templatesModified, id)
		if !ms.delete {
			continue
		}
		r.removeFile(path.Join(defaultRepo.PathTemplates, path.Base(ms.filename)))
	}
	return nil
}

func (r *repository) writeProjects() error {
	for _, val := range r.projects {
		if val.Id == "" {
			log.Fatalln("project has empty Id")
			continue
		}

		_, ok := r.projectsModified[val.Id]
		if !ok {
			continue
		}

		if val.FileName == "" {
			val.FileName = val.Id + ".yaml"
		}

		prj := proto.Clone(val).(*types.Project)
		if prj.Client != nil && prj.Client.Id != "" {
			prj.Client = &types.Client{Id: prj.Client.Id}
		}
		if prj.Account != nil && prj.Account.Id != "" {
			prj.Account = &types.Account{Id: prj.Account.Id}
		}

		buf, err := protojson.Marshal(prj)
		if err != nil {
			return err
		}

		buf, err = json2yaml(buf)
		if err != nil {
			return err
		}

		err = r.writeFile(path.Join(defaultRepo.PathProjects, path.Base(prj.FileName)), buf)
		if err != nil {
			return err
		}
	}

	for id, ms := range r.projectsModified {
		delete(r.projectsModified, id)
		if !ms.delete {
			continue
		}
		r.removeFile(path.Join(defaultRepo.PathProjects, path.Base(ms.filename)))
	}
	return nil
}

func (r *repository) writeInvoices() error {
	for _, val := range r.invoices {
		if val.Id == "" {
			log.Fatalln("invoice has empty Id")
			continue
		}

		_, ok := r.invoicesModified[val.Id]
		if !ok {
			continue
		}

		if val.FileName == "" {
			val.FileName = val.Id + ".yaml"
		}

		inv := proto.Clone(val).(*types.Invoice)
		if inv.Project != nil && inv.Project.Id != "" {
			inv.Project = &types.Project{Id: inv.Project.Id}
		}
		if inv.Account != nil && inv.Account.Id != "" {
			inv.Account = &types.Account{Id: inv.Account.Id}
		}

		buf, err := protojson.Marshal(inv)
		if err != nil {
			return err
		}

		buf, err = json2yaml(buf)
		if err != nil {
			return err
		}

		err = r.writeFile(path.Join(defaultRepo.PathInvoices, inv.Year, path.Base(inv.FileName)), buf)
		if err != nil {
			return err
		}
	}

	for id, ms := range r.invoicesModified {
		delete(r.invoicesModified, id)
		if !ms.delete {
			continue
		}
		r.removeFile(path.Join(defaultRepo.PathInvoices, path.Base(ms.filename)))
	}
	return nil
}
