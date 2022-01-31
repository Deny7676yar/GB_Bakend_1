На этом уроке мы создали: 
-server который посылает каждую секунду время,
подключившимся клиентам

-chart
Все вошедшие клиенты видят сообения которые посылает подключившийся клиент на сервер

1.Создали серверную часть
```
func main () {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil{
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil{
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}```
Горутина broadcaster хранит информацию о всех клиентах и прослушивает каналы событий и сообщений, 
используя мультиплексирование с помощью select
```
func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}
```
Горутина handleConn создает новый канал исходящих сообщений для своего клиента и 
объявляет широковещателю о поступлении этого клиента по каналу entуring.

```
func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	//realization client nickname
	who := conn.RemoteAddr().String()
	buffer := make([]byte, 100)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("could not read nickname from %s", who)
		return
	}
	nName := string(buffer)

	ch <- "You are " + nName
	messages <- nName + " has arrived"
	entering <- client{ch, nName}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- nName + ": " + input.Text()
	}
	leaving <- client{ch, nName}
	messages <- nName + " has left"
	conn.Close()
}```

2.Создали клиентскую часть
```
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin) // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
}
```