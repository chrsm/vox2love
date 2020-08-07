## vox2love

`vox2love` converts a MagicaVoxel `.vox` file into files loadable by Lovox.

Rather than loading an entire `.vox` at runtime, they can be preprocessed into
the format that Lovox expects.

Note that this currently only supports v150 of the `.vox` format, as specified
in `docs/file-format.txt`, which was pulled from [here][1].

Build the binary and run:

`vox2love -p path/to/vox -o path/to/destination`

You'll need to do some renaming for the files or individual frames for now.
I will probably update vox2love to create a directory structure later.

[1]: https://github.com/ephtracy/voxel-model
