# document
```
status code: 
200 // 서버에 파일이 존재합니다
403 // 서버에 접근이 차단되었습니다
404 // 서버에 파일이 없거나 해당 파일이 존재하지 않습니다
500 // 파라미터를 확인하거나 관리자에게 문의하세요
```

***
### 비트맵셋 다운로드
>/d/:id?nv=1

>성공/실패 판단기준 = status code
```
Param : {
    id int
}

QueryParam : {
    nv(NoVideo) bool // default bool
}
```
### note
> 서버가 가진 파일이 최신이 아닌경우 업데이트가 자동으로 진행됩니다.   
> 또한 서버의 다운로드 완료까지 기다릴 필요가 없습니다.

***

### 비트맵셋 검색 (미완성)
>/search
```
QueryParam:{
    p int page default 0
    
    sort string : {
        ranked_asc
        ranked_desc <= default
        favourites_asc
        favourites_desc
        plays_asc
        plays_desc
        updated_asc
        updated_desc
    }
    
    m int gameMode : [0,1,2,3] if this null select all
    s string : {
        ranked      1,2
        qualified   3
        loved       4
        pending     0
        wip         -1
        graveyard   -2
        any         4,3,2,1,0,-1,-2
        null        4,2,1
    }
    q string search strings	
    
    
}

response json body:{
    [
        mapsetId int
        ...
    ]
}
```
***
