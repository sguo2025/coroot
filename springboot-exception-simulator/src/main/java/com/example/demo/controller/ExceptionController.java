package com.example.demo.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.*;

@RestController
@RequestMapping("/simulate")
public class ExceptionController {

    @GetMapping("/runtime")
    public String runtimeError() {
        throw new RuntimeException("模拟运行时异常");
    }

    @GetMapping("/null")
    public String nullPointer() {
        String s = null;
        return s.toString();
    }

    @GetMapping("/delay")
    public String delay() throws InterruptedException {
        Thread.sleep(5000);
        return "延迟5秒返回";
    }

    @GetMapping("/oom")
    public String oom() {
        List<byte[]> list = new ArrayList<>();
        while (true) list.add(new byte[1024 * 1024]);
    }

    @GetMapping("/cpu")
    public String cpu() {
        long end = System.currentTimeMillis() + 5000;
        while (System.currentTimeMillis() < end) Math.pow(Math.random(), Math.random());
        return "CPU高占用结束";
    }

    @GetMapping("/random")
    public String random() throws InterruptedException {
        int r = new Random().nextInt(5);
        switch (r) {
            case 0:
                throw new RuntimeException("随机异常：RuntimeException");
            case 1:
                Thread.sleep(3000);
                return "随机延迟3秒返回";
            case 2:
                String s = null;
                return s.toString();
            default:
                return "正常返回";
        }
    }

    @GetMapping("/status/{code}")
    public ResponseEntity<String> status(@PathVariable int code) {
        return ResponseEntity.status(code).body("返回自定义状态码: " + code);
    }
}
