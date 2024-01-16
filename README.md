# Renum

A easy way to rename and re-numbering files in a directory.

## Usage (example)

Giving a folder with the following files:

```
ls -l /tmp/renum-test/
total 0
-rw-rw-r-- 1 alphayax alphayax 0 janv. 16 14:49 '[XXX-Fansub]_Xxx_Xxxxx_1086_[VOSTFR][FHD_1920x1080].xxx'
-rw-rw-r-- 1 alphayax alphayax 0 janv. 16 14:49 '[XXX-Fansub]_Xxx_Xxxxx_1087_[VOSTFR][FHD_1920x1080].xxx'
-rw-rw-r-- 1 alphayax alphayax 0 janv. 16 14:49 '[XXX-Fansub]_Xxx_Xxxxx_1088_[VOSTFR][FHD_1920x1080].xxx'
-rw-rw-r-- 1 alphayax alphayax 0 janv. 16 14:49 '[XXX-Fansub]_Xxx_Xxxxx_1089_[VOSTFR][FHD_1920x1080].xxx'
-rw-rw-r-- 1 alphayax alphayax 0 janv. 16 14:49 '[XXX-Fansub]_Xxx_Xxxxx_1090_[VOSTFR][FHD_1920x1080].xxx'
```

You can rename them with the following command:

```bash
renum --season 12 /tmp/renum-test/
```

And you will get:

```
2024/01/16 14:50:17 [Preview] [XXX-Fansub]_Xxx_Xxxxx_1086_[VOSTFR][FHD_1920x1080].xxx -> [XXX-Fansub]_Xxx_Xxxxx_S12E01_[VOSTFR][FHD_1920x1080].xxx
2024/01/16 14:50:17 [Preview] [XXX-Fansub]_Xxx_Xxxxx_1087_[VOSTFR][FHD_1920x1080].xxx -> [XXX-Fansub]_Xxx_Xxxxx_S12E02_[VOSTFR][FHD_1920x1080].xxx
2024/01/16 14:50:17 [Preview] [XXX-Fansub]_Xxx_Xxxxx_1088_[VOSTFR][FHD_1920x1080].xxx -> [XXX-Fansub]_Xxx_Xxxxx_S12E03_[VOSTFR][FHD_1920x1080].xxx
2024/01/16 14:50:17 [Preview] [XXX-Fansub]_Xxx_Xxxxx_1089_[VOSTFR][FHD_1920x1080].xxx -> [XXX-Fansub]_Xxx_Xxxxx_S12E04_[VOSTFR][FHD_1920x1080].xxx
2024/01/16 14:50:17 [Preview] [XXX-Fansub]_Xxx_Xxxxx_1090_[VOSTFR][FHD_1920x1080].xxx -> [XXX-Fansub]_Xxx_Xxxxx_S12E05_[VOSTFR][FHD_1920x1080].xxx
Do you want to continue the operation? (y/N): 
```
