# decent
Golang libs and executables for Decent DE1

## des

Server that exposes websockets of Decent events and serves console and debug pages. 

Console is available at: `http://localhost:8080`
Debug view, including simulating shots available at: `http://localhost:8080/debug`

```shell
% ./des
listening on :8080
```

To connect to a TCP socket outputting real DE1 events, use the `de1` switch:
```shell
% ./des -de1=192.168.1.20:19090
```

## dej

Simple utility to parse DE1 serial protocol from stdin and output in formatted JSON, one line per message.

```shell
% cat output.text | ./dej
{"type":"water_levels","water_level":11.472656,"water_fill_level":5}
{"type":"shot_sample","sample_time":14743,"group_pressure":0.0007324219,"group_flow":0,"mix_temp":90.14453,"head_temp":88.777084,"set_mix_temp":80,"set_head_temp":89,"set_group_pressure":0,"set_group_flow":3,"frame_number":3,"steam_temp":164}
...
```
