# Liar Of Turing

```
/chat-server
    /cmd
        /server
            main.go  // 채팅 서버 애플리케이션의 진입점
    /internal
        /websocket
            websocket.go  // 웹소켓 관리 (연결, 메시지 전송 등)
        /handler
            handler.go  // HTTP 요청 핸들러 (웹소켓 업그레이드 등)
        /room
            room.go    // 채팅방 관리 로직
    /pkg
        /logger
            logger.go  // 로깅 유틸리티
    README.md
```
