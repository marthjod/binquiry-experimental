package main

import (
	"log"
	"time"

	pb "github.com/marthjod/binquiry-experimental/word"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

var (
	word = []byte(`
<!-- <center> -->
<div class="page-header"><h2>kona <small>Kvenkynsnafnorð</small></h2></div>

<div class="row-fluid">
<div class="span6">
    <table cellpadding="6" class="table table-hover">
        <tr><th colspan="3" align="center">Eintala</th></tr>
        <tr>
        </tr>
        <td width="10%"></td>
        <td align="center" width="20%">&aacute;n&nbsp;greinis</td>
        <td align="center" width="20%">me&eth;&nbsp;greini</td>
        <tr>
            <td>Nf.</td>
            <td><span class="VO_beygingarmynd">kona</span></td>
            <td><span class="VO_beygingarmynd">konan</span></td>
        </tr>
        <tr>
            <td>Þf.</td>
            <td><span class="VO_beygingarmynd">konu</span></td>
            <td><span class="VO_beygingarmynd">konuna</span></td>
        </tr>
        <tr>
            <td>Þgf.</td>
            <td><span class="VO_beygingarmynd">konu</span></td>
            <td><span class="VO_beygingarmynd">konunni</span></td>
        </tr>
        <tr>
            <td>Ef.</td>
            <td><span class="VO_beygingarmynd">konu</span></td>
            <td><span class="VO_beygingarmynd">konunnar</span></td>
        </tr>
    </table>
</div>
<div class="span6">
    <table cellpadding="6" class="table table-hover">
        <tr><th colspan="3" align="center">Fleirtala</th></tr>
        <tr>
        </tr>
        <td width="10%"></td>
        <td align="center" width="20%">&aacute;n&nbsp;greinis</td>
        <td align="center" width="20%">me&eth;&nbsp;greini</td>
        <tr>
            <td>Nf.</td>
            <td><span class="VO_beygingarmynd">konur</span></td>
            <td><span class="VO_beygingarmynd">konurnar</span></td>
        </tr>
        <tr>
            <td>Þf.</td>
            <td><span class="VO_beygingarmynd">konur</span></td>
            <td><span class="VO_beygingarmynd">konurnar</span></td>
        </tr>
        <tr>
            <td>Þgf.</td>
            <td><span class="VO_beygingarmynd">konum</span></td>
            <td><span class="VO_beygingarmynd">konunum</span></td>
        </tr>
        <tr>
            <td>Ef.</td>
            <td><span class="VO_beygingarmynd">kvenna</span></td>
            <td><span class="VO_beygingarmynd">kvennanna</span></td>
        </tr>
    </table>
</div>
</div>


        <!-- </center> -->

`)
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDispatcherClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Dispatch(ctx, &pb.DispatchRequest{Word: word})
	if err != nil {
		log.Fatal("could not dispatch: ", err)
	}
	log.Printf("response: %v", r)
}
