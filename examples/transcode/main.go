package main

import (
	"github.com/strengine/Core/av"
	"github.com/strengine/Core/av/avutil"
	"github.com/strengine/Core/av/transcode"
	"github.com/strengine/Core/cgo/ffmpeg"
	"github.com/strengine/Core/format"
)

// need ffmpeg with libfdkaac installed

func init() {
	format.RegisterAll()
}

func main() {
	infile, _ := avutil.Open("speex.flv")

	findcodec := func(stream av.AudioCodecData, i int) (need bool, dec av.AudioDecoder, enc av.AudioEncoder, err error) {
		need = true
		dec, _ = ffmpeg.NewAudioDecoder(stream)
		enc, _ = ffmpeg.NewAudioEncoderByName("libfdk_aac")
		enc.SetSampleRate(stream.SampleRate())
		enc.SetChannelLayout(av.CH_STEREO)
		enc.SetBitrate(12000)
		enc.SetOption("profile", "HE-AACv2")
		return
	}

	trans := &transcode.Demuxer{
		Options: transcode.Options{
			FindAudioDecoderEncoder: findcodec,
		},
		Demuxer: infile,
	}

	outfile, _ := avutil.Create("out.ts")
	avutil.CopyFile(outfile, trans)

	outfile.Close()
	infile.Close()
	trans.Close()
}