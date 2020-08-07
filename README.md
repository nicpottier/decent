# decent
Golang libs and executables for Decent DE1

## dej

Simple utility to parse DE1 serial protocol from stdin and output in formatted JSON, one line per message.

```shell
% go build github.com/nicpottier/decent/cmd/dej
% cat output.text | ./dej
{"type":"water_levels","water_level":11.472656,"water_fill_level":5}
{"type":"shot_sample","sample_time":14743,"group_pressure":0.0007324219,"group_flow":0,"mix_temp":90.14453,"head_temp":88.777084,"set_mix_temp":80,"set_head_temp":89,"set_group_pressure":0,"set_group_flow":3,"frame_number":3,"steam_temp":164}
...
```