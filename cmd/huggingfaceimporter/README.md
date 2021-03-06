# Hugging Face Importer

## Import a Pre-Trained Model

spaGO allows you either to use a model in the inference phase or to train one from scratch, or fine-tune it. However,
training a language model (i.e. the transformer objective) to get competitive results can become prohibitive. This
applies in general, but even more so with spaGO as it does not currently use the GPU :scream:

Pre-trained (fine-tuned) transformer models exist for several languages and are publicly hosted on
the [Hugging Face models repository](https://huggingface.co/models).

Particularly, these exist for BERT, ELECTRA and BART the three types of transformers architectures currently supported
by spaGO.

## Build

Move into the top directory, and run the following command:

```console
GOARCH=amd64 go build -o huggingface-importer cmd/huggingfaceimporter/main.go 
```

## Usage

To import a pre-trained model, run the `huggingface-importer` indicating both the model name you'd like to import (
including organization), and a local directory where to store all your models.

Example:

```console
./huggingface-importer --model=deepset/bert-base-cased-squad2 --repo=~/.spago 
```

At the end of the process, you should see:

```console
Serializing model to "~/.spago/deepset/bert-base-cased-squad2/spago_model.bin"... ok
BERT has been converted successfully!
```

The directory `~/.spago/deepset/bert-base-cased-squad2` should contains the original Hugging Face files plus the files
generated by spaGO: `spago_model.bin` and `embeddings_storage`.

The Docker version can be run like this.

```console
docker run --rm -it -v ~/.spago:/tmp/spago spago:main huggingface-importer --model=deepset/bert-base-cased-squad2 --repo=/tmp/spago
```
