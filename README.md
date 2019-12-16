# ILLUMI-PACKET

パケットが光るLANケーブル

![illumi-packet](https://user-images.githubusercontent.com/29009733/70907987-8ab66000-204d-11ea-86e2-09a34d7c557a.jpg)

## 準備

### 必要なもの
|材料|量|諸注意|
|:-:|:-:|:-:|
|LEDテープ(WS281B)|1m|個別アドレス可能・フルカラーなもの．144LEDs/m 推奨．|
|Raspberry Pi|1台|動作確認: Raspberry Pi 3|
|LANケーブル|1m||
|ジャンパ線 オス-メス|3本||
|結束バンド|4本||

そのほかに、キーボード、ディスプレイ、HDMIケーブル、ルータ等はご用意ください。

### 環境設定
動作には以下の環境が必要です。
* golang
* libpcap
* SCons
* rpi_ws281x
* illumi-packet


#### golangのインストール

実行するプログラムはgo言語で記述されています。

https://golang.org/doc/install#install に従ってインストールします。

以下のようにターミナルで実行し、最後にバージョンが表示されることを確認してください。

```sh
$ version=1.13.4
$ wget https://storage.googleapis.com/golang/go${version}.linux-armv6l.tar.gz
$ sudo tar -C /usr/local -xzf go${version}.linux-armv6l.tar.gz

$ echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
$ . ~/.profile

$ go version
go version go1.13.4 linux/arm
```

#### libpcapのインストール
プログラムからパケットキャプチャを行うためにlibpcapをインストールします。

```sh
$ sudo apt-get install libpcap-dev
```

#### SConsのインストール
rpi_ws281xをビルドするためにSConsをインストールします。

```sh
$ sudo apt-get install scons
```

#### rpi_ws281xのインストール
LEDテープはrpi_ws281xというライブラリで操作します。

https://github.com/jgarff/rpi_ws281x に従ってインストールします。

```sh
$ git clone https://github.com/jgarff/rpi_ws281x.git
$ cd rpi_ws281x
$ scons

$ sudo cp -ai ./ws2811.h ./rpihw.h ./pwm.h /usr/local/include/
$ sudo cp -ai ./libws2811.a /usr/local/lib/
```

#### イルミパケットのソースコード
このリポジトリのソースコードを適当なディレクトリにダウンロードします。

```sh
$ git clone https://github.com/souring001/illumi-packet.git
$ cd illumi-packet
```

LEDの個数によって`illumi-packet.go`の以下の変数を適宜変更してください。

| LEDの個数 | count | speed | series |
| --------:| -----:| -----:| ------:|
|60 個/m   |    60 |      1 |     6 |
|144 個/m  |   144 |      4 |    12 |

## LANケーブルの作り方

1. LANケーブルにLEDテープを乗せて、結束バンドで固定する。
2. ジャンパワイヤ(オス側)を挿し込む
3. メス側を Raspberry Pi のGPIOの2(5V), 6(GND), 12(信号) に挿し込む
4. Raspberry Piとルータに接続する

![GPIO](https://user-images.githubusercontent.com/29009733/70908199-f7315f00-204d-11ea-9cb0-256967c7ca5e.png)

## ビルド方法
ソースコードを変更するたびにビルドをする必要があります。

```sh
$ go build illumi-packet.go
```

## 起動

```sh
$ sudo ./illumi-packet
```

### オプション

|オプション|内容|
|:-|:-|
|-h|オプションの説明|
|-debug |パケット情報等の詳細を表示する(デフォルトは`true`)|
|-device [string]|ネットワークインターフェースを設定(デフォルトは`eth0`)|
|-speed [int]|パケットの流れる速度を設定(デフォルトは`1`)|
|-narp|ARPを表示しない|
|-ntcp|TCPを表示しない|
|-nudp|UDPを表示しない|
|-ipaddr|IPアドレスをLEDに表示する|

例: IPアドレスをLEDに表示
```sh
$ sudo ./illumi-packet -ipaddr
```
結果:
![showipaddress](https://user-images.githubusercontent.com/29009733/70908359-5e4f1380-204e-11ea-9187-a2d385c9f300.JPG)
