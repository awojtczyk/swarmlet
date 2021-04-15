package docker

import (
	"context"
	"encoding/base64"
	"sort"

	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// SecretList return all secrets.
func SecretList(name string, pageIndex, pageSize int) (secrets []swarm.Secret, totalCount int, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		opts := types.SecretListOptions{}
		if name != "" {
			opts.Filters = filters.NewArgs()
			opts.Filters.Add("name", name)
		}
		secrets, err = cli.SecretList(ctx, opts)
		if err == nil {
			sort.Slice(secrets, func(i, j int) bool {
				return secrets[i].Spec.Name < secrets[j].Spec.Name
			})
			totalCount = len(secrets)
			start, end := misc.Page(totalCount, pageIndex, pageSize)
			secrets = secrets[start:end]
		}
		return
	})
	return
}

// SecretCreate create a secret.
func SecretCreate(info *model.ConfigCreateInfo) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		spec := swarm.SecretSpec{}
		spec.Name = info.Name
		if info.Base64 {
			spec.Data, err = base64.StdEncoding.DecodeString(info.Data)
			if err != nil {
				return
			}
		} else {
			spec.Data = []byte(info.Data)
		}
		spec.Labels = info.Labels.ToMap()
		if info.Template.Name != "" {
			spec.Templating = &swarm.Driver{
				Name:    info.Template.Name,
				Options: info.Template.Options,
			}
		}
		_, err = cli.SecretCreate(ctx, spec)
		return
	})
}

// SecretUpdate update a config.
func SecretUpdate(info *model.ConfigUpdateInfo) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var secret swarm.Secret
		secret, _, err = cli.SecretInspectWithRaw(ctx, info.ID)
		if err != nil {
			return err
		}

		spec := secret.Spec
		// only the Labels field can be updated on API 1.30
		//spec.Name = info.Name
		//spec.Data = []byte(info.Data)
		spec.Labels = info.Labels.ToMap()
		return cli.SecretUpdate(ctx, info.ID, version(info.Version), spec)
	})
}

// SecretInspect returns secret information with raw data.
func SecretInspect(id string) (secret swarm.Secret, raw []byte, err error) {
	var (
		ctx context.Context
		cli *client.Client
	)
	if ctx, cli, err = mgr.Client(); err == nil {
		secret, raw, err = cli.SecretInspectWithRaw(ctx, id)
	}
	return
}

// SecretRemove remove a secret.
func SecretRemove(id string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		return cli.SecretRemove(ctx, id)
	})
}
