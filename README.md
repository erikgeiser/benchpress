# Benchpress

Benchpress is a benchmark suite for Go compression libraries. It was developed
to get a feel for the performance of the Go implementations and the effect of
the parameters.

Currently the following compression algorithms are benchmarked:

- gzip
- xz
- zstd
- lzma
- lzma2

## Datasets

The benchmarks are run on any file present in the `testdata` folder. For
example, you can run the benchmarks on the [enwik8](https://cs.fit.edu/~mmahoney/compression/textdata.html) dataset:

```
wget http://mattmahoney.net/dc/enwik8.zip
unzip enwik8.zip -d testdata/
rm enwik8.zip
```
