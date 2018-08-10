// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package backups

import (
	"io"
	"os"
	"time"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/state/backups"
)

const (
	notset          = backups.FilenamePrefix + "<date>-<time>.tar.gz"
	downloadWarning = "downloading backup archives is recommended; " +
		"backups stored remotely are not guaranteed to be available."
)

const createDoc = `
This command requests that Juju creates a backup of its state and prints the
backup's unique ID.  You may provide a note to associate with the backup.

By default, the backup archive and associated metadata are downloaded 
without keeping a copy remotely on the controller.

Use --no-download to avoid getting a local copy of the backup downloaded 
at the end of the backup process.

Use --keep-copy option to store a copy of backup remotely on the controller.

Use --verbose to see extra information about backup.

To access remote backups stored on the controller, see 'juju download-backup'.

Examples:
    juju create-backup 
    juju create-backup --no-download
    juju create-backup --no-download --keep-copy=false // ignores --keep-copy
    juju create-backup --keep-copy
    juju create-backup --verbose

See also:
    backups
    download-backup
`

// NewCreateCommand returns a command used to create backups.
func NewCreateCommand() cmd.Command {
	return modelcmd.Wrap(&createCommand{})
}

// createCommand is the sub-command for creating a new backup.
type createCommand struct {
	CommandBase
	// NoDownload means the backups archive should not be downloaded.
	NoDownload bool
	// Filename is where the backup should be downloaded.
	Filename string
	// Notes is the custom message to associated with the new backup.
	Notes string
	// KeepCopy means the backup archive should be stored in the controller db.
	KeepCopy bool
	fs       *gnuflag.FlagSet
}

// Info implements Command.Info.
func (c *createCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "create-backup",
		Args:    "[<notes>]",
		Purpose: "Create a backup.",
		Doc:     createDoc,
	}
}

// SetFlags implements Command.SetFlags.
func (c *createCommand) SetFlags(f *gnuflag.FlagSet) {
	c.CommandBase.SetFlags(f)
	f.BoolVar(&c.NoDownload, noDownloadOption, false, "Do not download the archive, implies keep-copy")
	f.BoolVar(&c.KeepCopy, keepCopyOption, false, "Keep a copy of the archive on the controller")
	f.StringVar(&c.Filename, "filename", notset, "Download to this file")
	c.fs = f
}

var (
	noDownloadOption = "no-download"
	keepCopyOption   = "keep-copy"
)

// Init implements Command.Init.
func (c *createCommand) Init(args []string) error {
	noDownloadSet, keepCopySet := false, false
	c.fs.Visit(func(blah *gnuflag.Flag) {
		if blah.Name == noDownloadOption {
			noDownloadSet = true
		}
		if blah.Name == keepCopyOption {
			keepCopySet = true
		}
	})
	if keepCopySet && noDownloadSet {
		return errors.Errorf("cannot mix --no-download and --keep-copy")
	}
	notes, err := cmd.ZeroOrOneArgs(args)
	if err != nil {
		return err
	}
	c.Notes = notes

	if c.Filename != notset && c.NoDownload {
		return errors.Errorf("cannot mix --no-download and --filename")
	}

	if c.Filename == "" {
		return errors.Errorf("missing filename")
	}
	return nil
}

// Run implements Command.Run.
func (c *createCommand) Run(ctx *cmd.Context) error {
	if c.Log != nil {
		if err := c.Log.Start(ctx); err != nil {
			return err
		}
	}

	client, apiVersion, err := c.NewGetAPI()
	if err != nil {
		return errors.Trace(err)
	}
	defer client.Close()

	if apiVersion < 2 {
		if c.KeepCopy {
			return errors.New("--keep-copy is not supported by this controller")
		}
		// for API v1, keepCopy is the default and only choice, so set it here
		c.KeepCopy = true
	}

	if c.NoDownload {
		ctx.Warningf(downloadWarning)
		c.KeepCopy = true
	}

	metadataResult, copyFrom, err := c.create(client, apiVersion)
	if err != nil {
		return errors.Trace(err)
	}

	// TODO: (hml) 2018-04-25
	// fix to dump the metadata when --verbose used
	if c.Log != nil && !c.Log.Quiet {
		c.dumpMetadata(ctx, metadataResult)
	}

	if c.KeepCopy {
		ctx.Infof("Remote backup stored on the controller as %v.", metadataResult.ID)
	} else {
		ctx.Infof("Remote backup was not created.")
	}

	// Handle download.
	if !c.NoDownload {
		filename := c.decideFilename(ctx, c.Filename, metadataResult.Started)
		if err := c.download(ctx, client, copyFrom, filename); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (c *createCommand) decideFilename(ctx *cmd.Context, filename string, timestamp time.Time) string {
	if filename != notset {
		return filename
	}
	// Downloading but no filename given, so generate one.
	return timestamp.Format(backups.FilenameTemplate)
}

func (c *createCommand) download(ctx *cmd.Context, client APIClient, copyFrom string, archiveFilename string) error {
	resultArchive, err := client.Download(copyFrom)
	if err != nil {
		return errors.Trace(err)
	}
	defer resultArchive.Close()

	archive, err := os.Create(archiveFilename)
	if err != nil {
		return errors.Annotatef(err, "while creating local archive file %v", archiveFilename)
	}
	defer archive.Close()

	_, err = io.Copy(archive, resultArchive)
	if err != nil {
		return errors.Annotatef(err, "while copying to local archive file %v", archiveFilename)
	}
	ctx.Infof("Downloaded to %v.", archiveFilename)
	return nil
}

func (c *createCommand) create(client APIClient, apiVersion int) (*params.BackupsMetadataResult, string, error) {
	result, err := client.Create(c.Notes, c.KeepCopy, c.NoDownload)
	if err != nil {
		return nil, "", errors.Trace(err)
	}
	copyFrom := result.ID

	if apiVersion >= 2 {
		copyFrom = result.Filename
	}

	return result, copyFrom, err
}
