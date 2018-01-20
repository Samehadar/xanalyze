SET DEFINE OFF;
CREATE TABLE D_CITY_CODES
(
  CC_CODE          VARCHAR2(3 BYTE),
  DESCRIPTION      VARCHAR2(200 BYTE),
  COUNTRYCODE      VARCHAR2(2 BYTE),
  OUTBOUND         VARCHAR2(1 BYTE),
  DESTINATION      VARCHAR2(1 BYTE),
  CITY             VARCHAR2(1 BYTE),
  COUNTRY          VARCHAR2(1 BYTE),
  GOUPING_LEVEL1   VARCHAR2(250 BYTE),
  ADDITIONAL_DATA  VARCHAR2(1000 BYTE)
)
TABLESPACE DTW_ADV_TABLES
RESULT_CACHE (MODE DEFAULT)
PCTUSED    0
PCTFREE    10
INITRANS   1
MAXTRANS   255
STORAGE    (
            INITIAL          80K
            NEXT             1M
            MINEXTENTS       1
            MAXEXTENTS       UNLIMITED
            PCTINCREASE      0
            BUFFER_POOL      DEFAULT
            FLASH_CACHE      DEFAULT
            CELL_FLASH_CACHE DEFAULT
           )
LOGGING 
NOCOMPRESS 
NOCACHE
NOPARALLEL
MONITORING;


CREATE UNIQUE INDEX D_CITY_CODES_PK ON D_CITY_CODES
(CC_CODE)
LOGGING
TABLESPACE DTW_ADV_TABLES
PCTFREE    10
INITRANS   2
MAXTRANS   255
STORAGE    (
            INITIAL          80K
            NEXT             1M
            MINEXTENTS       1
            MAXEXTENTS       UNLIMITED
            PCTINCREASE      0
            BUFFER_POOL      DEFAULT
            FLASH_CACHE      DEFAULT
            CELL_FLASH_CACHE DEFAULT
           )
NOPARALLEL;


ALTER TABLE D_CITY_CODES ADD (
  CONSTRAINT D_CITY_CODES_PK
  PRIMARY KEY
  (CC_CODE)
  USING INDEX D_CITY_CODES_PK
  ENABLE VALIDATE);