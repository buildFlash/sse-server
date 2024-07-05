package main

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"log"
	"time"
)

var index = []byte(`<!DOCTYPE html>
<html>
<body>
<h1>SSE Messages</h1>
<div id="result"></div>
<script>
if(typeof(EventSource) !== "undefined") {
  var source = new EventSource("http://127.0.0.1:3000/sse");
  source.onmessage = function(event) {
    document.getElementById("result").innerHTML += event.data + "<br>";
  };
} else {
  document.getElementById("result").innerHTML = "Sorry, your browser does not support server-sent events...";
}
</script>
</body>
</html>
`)

func main() {
	// Fiber instance
	app := fiber.New()

	// CORS for external resources
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Cache-Control"},
	}))

	app.Get("/", func(ctx fiber.Ctx) error {
		ctx.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		return ctx.Status(fiber.StatusOK).Send(index)
	})

	app.Get("/sse", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			fmt.Println("WRITER")
			var i int
			for {
				i++
				msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
				_, err := fmt.Fprintf(w, "data: Message: %s\n\n", msg)
				if err != nil {
					log.Printf("Error: %s", err)
					return
				}
				fmt.Println(msg)

				w.Flush()
				time.Sleep(2 * time.Second)
			}
		})

		return nil
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
