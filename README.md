# pinyin-of

Find the pinyin of a Chinese character. Contains the following:

* The original source data (obtained from https://download.fcitx-im.org/data/dict.utf8-20170423.tar.xz), `dict.utf8`;

* A tool (`dict-to-pinyin-map`) to convert the source data to a file more suitable for the purpose of this tool, `pinyin-map`;

* A tool (`pinyin-of`) that uses `pinyin-map` to convert Chinese characters to pinyin.
