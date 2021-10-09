# Docker Compose 테스트

## `proxy_net_test.sh`

`share/pnet`은 TCP 연결 중 local port를 변경하는 iptables SNAT rule을 추가하는 작업을 한다.
iptables를 쓸 수 있는 실행 환경과 외부 서버가 필요하다.

### 컨테이너 구성

Compose로 다음 세 개의 역할을 하는 컨테이너들을 올린다.

* Agent Client

  `share/pnet`패키지를 테스트한다.

* router

  Agent Client의 네트워크를 SNAT한다. (이 테스트에서 반드시 필요하진 않으나 나중에 다른 테스트에 필요하기 때문에 미리 준비했다.)

* echo server

  TCP 연결을 검증하기 위한 간단한 `ncat` 서버를 실행한다.

세부 동작 흐름은 코드를 참조한다.
