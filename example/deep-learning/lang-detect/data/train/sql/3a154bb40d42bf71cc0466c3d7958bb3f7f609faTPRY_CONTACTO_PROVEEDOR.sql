
  CREATE TABLE "REFERENCIA"."TPRY_CONTACTO_PROVEEDOR" 
   (	"IDCONTACTO" NUMBER NOT NULL ENABLE, 
	"IDPROVEEDOR" NUMBER NOT NULL ENABLE, 
	"NOMBRECONTACTO" VARCHAR2(50 BYTE) NOT NULL ENABLE, 
	"APPATERNO" VARCHAR2(50 BYTE) NOT NULL ENABLE, 
	"APMATERNO" VARCHAR2(50 BYTE) NOT NULL ENABLE, 
	"EMAIL1" VARCHAR2(50 BYTE) NOT NULL ENABLE, 
	"TELEFONO" VARCHAR2(10 BYTE) NOT NULL ENABLE, 
	"MOVIL1" VARCHAR2(13 BYTE), 
	"MOVIL2" VARCHAR2(13 BYTE), 
	"PRINCIPAL" VARCHAR2(1 BYTE), 
	 CONSTRAINT "TPRY_CONTACTOSPROVEEDOR_PK" PRIMARY KEY ("IDCONTACTO")
  USING INDEX PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS 
  STORAGE(INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645
  PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1
  BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT)
  TABLESPACE "DATA"  ENABLE, 
	 CONSTRAINT "FK001_TPRY_CONTACTO_PROVEEDOR" FOREIGN KEY ("IDPROVEEDOR")
	  REFERENCES "REFERENCIA"."TPRY_PROVEEDORES" ("IDPROVEEDOR") ENABLE
   ) SEGMENT CREATION IMMEDIATE 
  PCTFREE 10 PCTUSED 40 INITRANS 1 MAXTRANS 255 
 NOCOMPRESS LOGGING
  STORAGE(INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645
  PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1
  BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT)
  TABLESPACE "DATA" ;