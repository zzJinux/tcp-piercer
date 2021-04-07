# TCP Piercer
TCP Piercer는 NAT/Firewall 우회를 도울 수 있는 client-server 프로그램입니다.

일반적으로 라우터의 NAT/Firewal는 outbound TCP request는 허용되고 inbound TCP request가 차단하도록 설정됩니다.

<img src="./docs/_images/given-situation.png" alt="given situation" width=600>

TCP Piercer는 라우터 설정에 접근할 필요 없이 내부 네트워크의 호스트가 inbound TCP request도 받을 수 있게 하며
터널링 없이 TCP connection이 유지 됩니다.

<img src="./docs/_images/desired-result.png" alt="desired result" width=800>

후술할 한계점 때문에 터널링 프로토콜을 완전히 대체할 수 없습니다. TCP를 심층적으로 이해하는 것에 의의가 있는 개념증명 구현입니다.

## 전제 조건
- 라우터의 NAT translation behavior는 [endpoint-independent][rfc4787, section 4.1]
- 라우터와 인터넷 사이의 _모든 패킷 흐름을 통제할 수 있는_ 장비가 필요함 (터널링 프로토콜과 비교되는 한계점)

## 구현
- [네트워크 내부 흐름](./docs/network-flow.md)
- [TUN/TAP device][Wikipedia, TUN/TAP]

## 시뮬레이션
- Docker
- `debian:buster` containers


[rfc4787, section 4.1]: https://tools.ietf.org/html/rfc4787#section-4.1
[Wikipedia, TUN/TAP]: https://en.wikipedia.org/wiki/TUN/TAP

