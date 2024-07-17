# 20240222

0.0.22

## 버그

### LogZipProxyThread
* queue에 더이상 항목이 없는 경우 . 기존은 버퍼가 다 찰때까지 전송안되는 오류 수정.
* queue timeout (기본 2초) 인 경우 무조건 발송 

### dateutil 
* DateFormat 년도 오류 수정. 

     
--------------------------------------------------------------------------------

# 20240417

0.0.23

## 업데이트 

### RequestDoublieQueue
* GetTimeout 함수 추가 

### util.exception  
* CustomException 추가 
* Go 에서 error interface 를 추가한 오류 관련 정보

     
--------------------------------------------------------------------------------

# 20240516

0.0.24

## 버그 

### util.exception 
* message 입력 오류 수정


## 업데이트 

### create pack 에 LogSinkZipPack 추가

### ZipSendPtoxyTrhead compressutil 로 압축 함수 변경


--------------------------------------------------------------------------------

# 20240718

0.0.25

## 업데이트 

### StatGeneralPack Type 추가, STAT_GETNERAL_1 타입 추가. 
### HashUtil 에 V2 함수 이름만 변경해서 추가 (Hash64V2, Hash64StrV2)
     
--------------------------------------------------------------------------------