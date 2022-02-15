package org.proj.api.ws;

import org.java_websocket.WebSocket;
import org.java_websocket.client.WebSocketClient;
import org.java_websocket.handshake.ClientHandshake;
import org.java_websocket.handshake.ServerHandshake;
import org.java_websocket.server.WebSocketServer;

import java.net.InetSocketAddress;
import java.net.URI;

public class RegisterListener extends WebSocketServer {
    public RegisterListener (int port) {
        super(new InetSocketAddress(port));
    }

    @Override
    public void onOpen (WebSocket webSocket, ClientHandshake clientHandshake) {
        webSocket.send("Connected to " + webSocket.getRemoteSocketAddress());
    }

    @Override
    public void onClose (WebSocket webSocket, int i, String s, boolean b) {
        System.out.println("Closing server");
    }

    @Override
    public void onMessage (WebSocket webSocket, String s) {
        System.out.println("Server: " + s);
        webSocket.send("Hola bb ;)");
    }

    @Override
    public void onError (WebSocket webSocket, Exception e) {
        System.err.println("Server error: "+e.getLocalizedMessage());
    }

    @Override
    public void onStart() {
        System.out.println("Server started");
    }

    public static void main (String... args) {
        var server = new RegisterListener(1234);
        server.start();
    }
}

class TestingClient extends WebSocketClient {
    public TestingClient (String path) {
        super(URI.create("ws://" + path));
    }

    @Override
    public void onOpen (ServerHandshake serverHandshake) {
        System.out.println("Client opened");
    }

    @Override
    public void onMessage (String s) {
        System.out.println("Client: " + s);
    }

    @Override
    public void onClose (int i, String s, boolean b) {
        System.out.println("Closing client");
    }

    @Override
    public void onError (Exception e) {
        System.err.println("Client error: "+e.getLocalizedMessage());
    }
}