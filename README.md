# TCP Piercer
Router의 NAT에 인위적으로 NAT session을 생성을 유도하여 TCP 연결 수립을 정상적으로 완료시키고,
동시에 해당 연결의 성능 하락을 최소화할 수 있는 기법을 제시합니다.

> TCP Piercer는 TCP 동작을 이용하여 구현할 수 있는 새로운 기능한 개념 증명의 성격이 강하며
> production level 구현이 (당장은) 고려되지 않습니다.

---
TCP Piercer는 기존 라우터의 NAT/Firewall에 차단되지 않게 안팎의 호스트 간 TCP 연결 수립을 중개합니다.

일반적으로 라우터의 NAT/Firewall은 outbound TCP request는 허용하고 inbound TCP request가 차단하도록 설정됩니다.
<br>
<br>
<img src="./docs/_images/intro1.png" alt="intro1 - connection refused" width="600">

TCP Piercer는 라우터 설정에 접근할 필요 없이 내부 네트워크의 호스트가 inbound TCP request를 수신하는 것이 가능하게 하며,
이 때 수립된 TCP 연결은 [터널링][Wikipedia, Tunneling protocol] 없이 동작하기 때문에 오버헤드가 없습니다.
<br>
<br>
<img src="./docs/_images/intro2.png" alt="intro2 - we got a connection">


## 설계
자세한 동작 원리는 [이 곳](./docs/network-flow.md)에 나와있습니다.
- 라우터의 NAT translation behavior는 [endpoint-independent][rfc4787, section 4.1]라고 가정합니다.
- Agent Server는 라우터와 인터넷 사이의 _모든 패킷 흐름을 통제할 수_ 있다고 가정합니다.
- Agent Server가 수행하는 NAT, ISN 조작은 [tun device][Wikipedia, TUN/TAP]를 통해 이루어집니다.


## 한계
- Agent Server는 라우터와 인터넷 사이의 _모든 패킷 흐름을 통제할 수_ 있어야 합니다. 따라서 모든 상황에서
터널링 프로토콜을 대체할 수 있는 것은 아닙니다.

## 계획
### 시뮬레이션
TCP Piercer는 실제 네트워크 장비에서 동작하는 것이 이상적이나 본 프로젝트가 개념 증명에 초점을 맞췄기 때문에
빠르고 효율적인 실증을 위해 여러 개의 도커 컨테이너로 시뮬레이션할 수 있는 도구도 포함할 것입니다.

vpn 같은 터널링 프로토콜의 2-way connection을 없애기위해 router의 NAT에 fake mapping을
만들어 vpn 같은 기능을 제공하는 기법을 소개한다.


[rfc4787, section 4.1]: https://tools.ietf.org/html/rfc4787#section-4.1
[Wikipedia, TUN/TAP]: https://en.wikipedia.org/wiki/TUN/TAP
[Wikipedia, Tunneling protocol]: https://en.wikipedia.org/wiki/Tunneling_protocol
