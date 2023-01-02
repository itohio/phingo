package repository

import (
	"errors"
	"fmt"
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
	_ = r.fs.Remove(pth + ".rm")
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
	for _, val := range r.templates {
		if val.Id == "" {
			log.Fatalln("template has empty Id")
			continue
		}

		_, ok := r.templatesModified[val.Id]
		if !ok {
			continue
		}

		if val.FileName == "" {
			val.FileName = val.Id + ".md"
		}
		val.FileName = types.SanitizePath(val.FileName)
		err := r.writeFile(path.Join(defaultRepo.PathTemplates, path.Base(val.FileName)), val.Text)
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
		val.FileName = types.SanitizePath(val.FileName)

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
			return errors.New("id must be already set")
		}
		if val.Account == nil {
			return errors.New("account must be set")
		}
		if val.Client == nil {
			return errors.New("client must be set")
		}

		_, ok := r.invoicesModified[val.Id]
		if !ok {
			continue
		}

		val = proto.Clone(val).(*types.Invoice)
		val.Project = &types.Project{
			Id: val.Project.Id,
		}

		if val.FileName == "" {
			val.FileName = val.Code + ".yaml"
		}
		val.FileName = types.SanitizePath(val.FileName)
		buf, err := protojson.Marshal(val)
		if err != nil {
			return err
		}

		buf, err = json2yaml(buf)
		if err != nil {
			return err
		}

		err = r.writeFile(path.Join(defaultRepo.PathInvoices, fmt.Sprint(val.Year()), path.Base(val.FileName)), buf)
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
