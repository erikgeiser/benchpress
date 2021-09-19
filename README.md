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

## Running the benchmarks

The benchmarks can be run with the following command:
```
go test -bench="."
```

For nicer output, [prettybench](https://github.com/cespare/prettybench) can be used.

## Results

Example results on an `Intel i7 4770k` with `16 GB RAM`:

**Speed:**
```
benchmark                                                           iter        time/iter
---------                                                           ----        ---------
BenchmarkGzip/enwik8_compress_gzip_level_"best_speed"-8                1    1255.28 ms/op
BenchmarkGzip/enwik8_decompress_gzip_level_"best_speed"-8              1    1002.77 ms/op
BenchmarkGzip/enwik8_compress_gzip_level_"best_compression"-8          1    4348.42 ms/op
BenchmarkGzip/enwik8_decompress_gzip_level_"best_compression"-8        2     831.58 ms/op
BenchmarkGzip/enwik8_compress_gzip_level_"default"-8                   1    3911.91 ms/op
BenchmarkGzip/enwik8_decompress_gzip_level_"default"-8                 2     830.93 ms/op
BenchmarkZstd/enwik8_compress_zstd_level_"fastest"-8                   3     482.21 ms/op
BenchmarkZstd/enwik8_decompress_zstd_level_"fastest"-8                 2     537.83 ms/op
BenchmarkZstd/enwik8_compress_zstd_level_"default"-8                   2     654.13 ms/op
BenchmarkZstd/enwik8_decompress_zstd_level_"default"-8                 1    1381.82 ms/op
BenchmarkZstd/enwik8_compress_zstd_level_"better_compression"-8        1    1936.95 ms/op
BenchmarkZstd/enwik8_decompress_zstd_level_"better_compression"-8      1    1035.87 ms/op
BenchmarkZstd/enwik8_compress_zstd_level_"best_compression"-8          1    8443.98 ms/op
BenchmarkZstd/enwik8_decompress_zstd_level_"best_compression"-8        1    1731.81 ms/op
BenchmarkXZ/enwik8_compress_xz-8                                       1   15279.04 ms/op
BenchmarkXZ/enwik8_decompress_xz-8                                     1    4795.88 ms/op
BenchmarkLZMA/enwik8_compress_lzma-8                                   1   15561.28 ms/op
BenchmarkLZMA/enwik8_decompress_lzma-8                                 1    4048.16 ms/op
BenchmarkLZMA2/enwik8_compress_lzma2-8                                 1   15283.31 ms/op
BenchmarkLZMA2/enwik8_decompress_lzma2-8                               1    4656.62 ms/op
ok      benchpress      162.834s
```

**Size:**
```
enwik8: initial size: 95 MB

enwik8: gzip level "best speed":         40.81 MB (42.79%)
enwik8: gzip level "best compression":   34.95 MB (36.65%)
enwik8: gzip level "default":            34.99 MB (36.69%)
enwik8: zstd level "default":            34.37 MB (36.04%)
enwik8: zstd level "better compression": 31.83 MB (33.37%)
enwik8: zstd level "best compression":   28.59 MB (29.98%)    
enwik8: xz:                              31.78 MB (33.32%)
enwik8: lzma:                            31.77 MB (33.31%)
enwik8: lzma2:                           31.78 MB (33.32%)
```