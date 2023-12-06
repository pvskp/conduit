![plot](./assets/logo.svg)

# Conduit
Conduit is a powerful and versatile command-line application for HTTP routing and load balancing. Designed to be lightweight, easy to use, and highly configurable, Conduit is ideal for developers and network administrators who need an effective solution for traffic handling, load distribution, and basic network request manipulation.

### Features
- HTTP Proxy: Forward HTTP requests to specific backend servers.
- Load Balancing: Distribute incoming traffic across multiple servers using algorithms such as Round Robin, Least Connections, and Least Time.

## Installation
Follow these steps to install Conduit on your system.

### Prerequisites
- Go 1.15 or higher.

### Installation Instructions
1. Clone the repository:
```bash
$ git clone https://github.com/your-username/conduit.git
```
2. Navigate to the project directory:
```bash
$ cd conduit
```

3. Build the project:
```bash
$ go build
```

## Usage

Here are some examples of how you can use Conduit:

### Start a HTTP Proxy
```bash
$ conduit proxy -p 8080
```

### Start a load balancer with Round Robin algorithm
```bash
$ conduit loadBalancer -a roundRobin -H http://server1.example.com,http://server2.example.com -p 8080
```

## Contributing
Contributions are always welcome! If you have a suggestion that would improve this project, follow these steps:

1. Fork the project repository.
2. Create a new branch (git checkout -b feature/AmazingFeature).
3. Make your changes.
4. Commit your changes (git commit -m 'Add some AmazingFeature').
5. Push to the branch (git push origin feature/AmazingFeature).
6. Open a pull request.

## License
Distributed under the MIT License. See LICENSE for more information.
