![alt text](https://raw.githubusercontent.com/nlpodyssey/spago/main/assets/spago_logo.png)

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/nlpodyssey/spago/.svg)](https://pkg.go.dev/github.com/nlpodyssey/spago/)
![Go](https://github.com/nlpodyssey/spago/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nlpodyssey/spago)](https://goreportcard.com/report/github.com/nlpodyssey/spago)
[![Maintainability](https://api.codeclimate.com/v1/badges/be7350d3eb1a6a8aa503/maintainability)](https://codeclimate.com/github/nlpodyssey/spago/maintainability)
[![codecov](https://codecov.io/gh/nlpodyssey/spago/branch/main/graph/badge.svg)](https://codecov.io/gh/nlpodyssey/spago)
[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)
![Unstable](https://raw.githubusercontent.com/nlpodyssey/spago/main/assets/stability-unstable-yellow.svg)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

<p align="center"><i>If you like the project, please ★ star this repository to show your support! 🤩</i></p>

A **Machine Learning** library written in pure Go designed to support relevant neural architectures in **Natural
Language Processing**.

spaGO is self-contained, in that it uses its own lightweight *computational graph* framework for both training and
inference, easy to understand from start to finish.

## Usage

Requirements:

* [Go 1.15](https://golang.org/dl/)
* [Go modules](https://blog.golang.org/using-go-modules)

Clone this repo or get the library:

```console
go get -u github.com/nlpodyssey/spago
```

spaGO supports two main use cases, which are explained more in detail in the following.

### CLI mode

Several programs can be leveraged to tour the current NLP capabilities in spaGO. A list of the demos now follows.

* [Named Entities Recognition](https://github.com/nlpodyssey/spago/tree/main/cmd/ner)
* [Hugging Face Importer](https://github.com/nlpodyssey/spago/tree/main/cmd/huggingfaceimporter)
* [Question Answering](https://github.com/nlpodyssey/spago/tree/main/cmd/bert#question-answering-task)
* [Masked Language Model](https://github.com/nlpodyssey/spago/tree/main/cmd/bert#masked-language-model)

The Docker image can be built like this.

```console
docker build -t spago:main . -f Dockerfile
```

### Deploy On Okteto Cloud With One Click 

[![Develop on Okteto](https://okteto.com/develop-okteto.svg)](https://cloud.okteto.com/deploy?repository=https://github.com/nlpodyssey/spago/&branch=okteto-stack-v0.4.1)



### Library mode

You can access the core functionality of spaGO, i.e. optimizing mathematical expressions by back-propagating gradients
through a computational graph, in your own code by using spaGO in library mode.

At a high level, it comprises four main modules:

1. Matrix
2. Graph
3. Model
4. Optimizer

To get started, look at the implementation of built-in neural models, such as
the [LSTM](https://github.com/nlpodyssey/spago/blob/main/pkg/ml/nn/recurrent/lstm/lstm.go). Don't be afraid, it is
straightforward Go code. The idea is that you could have written spaGO :)

You may find a [Feature Source Tree](https://github.com/nlpodyssey/spago/blob/main/FEATURE_SOURCE_TREE.md) useful for a
quick overview of the library's package organization.

There is also a [repo](https://github.com/nlpodyssey/spago-examples) with handy examples, such as MNIST classification.

## Features

#### Automatic differentiation

- You write the *forward()*, it does all *backward()* derivatives for you:
    - Define-by-Run (default, just like PyTorch does)
    - Define-and-Run (similar to the static graph of TensorFlow)

#### Optimization methods

- Gradient descent:
    - Adam, RAdam, RMS-Prop, AdaGrad, SGD
- Differential Evolution

#### Neural networks
-   Feed-forward models (Linear, Highway, Convolution, ...)
-   Recurrent models (LSTM, GRU, BiLSTM...)
-   Attention mechanisms (Self-Attention, Multi-Head Attention, ...)
-   Recursive auto-encoders

#### Natural Language Processing
-   Memory-efficient Word Embeddings (with [badger](https://github.com/dgraph-io/badger) key–value store)
-   Character Language Models
-   Recurrent Sequence Labeler with CRF on top (e.g. Named Entities Recognition)
-   Transformer models (BERT-like)
    -   Masked language model
    -   Next sentence prediction
    -   Tokens Classification
    -   Text Classification (e.g. Sentiment Analysis)
    -   Question Answering
    -   Textual Entailment
    -   Text Similarity

#### Compatible with pre-trained state-of-the-art neural models:

- Hugging Face [Transformers](https://github.com/huggingface/transformers)
- [Flair](https://github.com/flairNLP/flair) sequence labeler architecture

## Current Status

We're not at a v1.0.0 yet, so spaGO is currently work-in-progress.

However, it has been running smoothly for a quite a few months now in a system that analyzes thousands of news items a
day!

Besides, it's pretty easy to get your hands on through, so you might want to use it in your real applications.

Early adopters may make use of it for production use today as long as they understand and accept that spaGO is not fully
tested and that APIs might change.

## Known Limits

Sadly, at the moment, spaGO is not GPU friendly by design.

## Contributing

We're glad you're thinking about contributing to spaGO! If you think something is missing or could be improved, please
open issues and pull requests. If you'd like to help this project grow, we'd love to have you!

To start contributing, check
the [Contributing Guidelines](https://github.com/nlpodyssey/spago/blob/main/CONTRIBUTING.md).

## Contact

We encourage you to write an issue. This would help the community grow.

If you really want to write to us privately, please email [Matteo Grella](mailto:matteogrella@gmail.com) with your
questions or comments.

## Links

* [spaGO: Self-contained ML/NLP Library in Go](https://www.slideshare.net/MatteoGrella/spago-a-selfcontained-ml-nlp-library-in-go) ([video of GoWayFest 4.0](https://www.youtube.com/watch?v=wE3CQU4G2fk))

## Acknowledgments

spaGO is a personal project that is part of the open-source [NLP Odyssey](https://github.com/nlpodyssey) initiative
initiated by members of the EXOP team. I would therefore like to thank [EXOP GmbH](https://www.exop-group.com/en/) here,
which is providing full support for development by promoting the project and giving it increasing importance.

## Sponsors

We appreciate contributions of all kinds. We especially want to thank spaGO fiscal sponsors who contribute to ongoing
project maintenance.

![Faire.ai logo](https://raw.githubusercontent.com/nlpodyssey/spago/main/assets/sponsors/faire_ai_logo.png)

> Our aim is simplifying people's life by making lending easy, fast and accessible, leveraging Open Banking and Artificial Intelligence.
> https://www.faire.ai/

![Hype logo](https://raw.githubusercontent.com/nlpodyssey/spago/main/assets/sponsors/hype_design_logo.png)

> We work on Artificial Intelligence based hardware and software systems, declining them in areas such as Energy Management, Personal Safety, E-Health and Sports equipment.
> https://hype-design.it/

![BoxxApps logo](https://raw.githubusercontent.com/nlpodyssey/spago/main/assets/sponsors/boxxapps_logo.png)

> Professional services in the IT sector for Local Administrations, Enterprises and Local Authorities.
> https://www.boxxapps.com/

See our [Open Collective](https://opencollective.com/nlpodyssey/contribute) page if you too are interested in becoming a sponsor.
