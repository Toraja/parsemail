package main

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

var to = []string{
	"bob@example.com",
	"cora@example.com",
	"david@example.com",
	"ema@example.com",
	"fox@example.com",
	"giraffe@example.com",
	"hen@example.com",
	"ice@example.com",
	"jack@example.com",
	"kate@example.com",
	"linda@example.com",
	"monica@example.com",
	"nancy@example.com",
	"owen@example.com",
	"patric@example.com",
	"queen@example.com",
	"richard@example.com",
	"susan@example.com",
	"tom@example.com",
	"uva@example.com",
	"vicky@example.com",
	"wane@example.com",
	"xi@example.com",
	"yoda@example.com",
	"zebra@example.com",
}

func TestSend1(t *testing.T) {
	if err := send1(to...); err != nil {
		t.Fatal(err)
	}
}

func TestSend2(t *testing.T) {
	if err := send2(to...); err != nil {
		t.Fatal(err)
	}
}

func TestSend3(t *testing.T) {
	if err := send3(to...); err != nil {
		t.Fatal(err)
	}
}

func TestDialAndSend1(t *testing.T) {
	if err := dialAndSend1(to...); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSend1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send1(to...)
	}
}

func BenchmarkSend2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := send2(to...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSend3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := send3(to...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDialAndSend1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := dialAndSend1(to...); err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkDialAndSend2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := dialAndSend2(to...); err != nil {
			b.Fatal(err)
		}
	}
}

func TestEml(t *testing.T) {
	msg, err := parseMail("output/singlemail.eml")
	if err != nil {
		t.Fatal(err)
	}
	eb, err := parseMultipart(msg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("------------------------")
	fmt.Printf("%+v\n", msg)
	fmt.Println("------------------------")
	fmt.Printf("Text:\n%s\nAttachment:\n%s\n", eb.Text, eb.Attachment)

	attach := "apple,banana,cherry,durian,elderberry,fig,grape\narmadillo,bear,cat,dog,elephant,fox,gopher\nbear,gin,grappa,mead,rum,tequilawine,whisky\n"
	assert.Equal(t, attach, string(eb.Attachment))
}
