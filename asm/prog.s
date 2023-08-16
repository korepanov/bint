.data
starSymbol:
.ascii "*"
endSymbol:
.ascii ";"
pageSize:
.quad 4096
varNameSize:
.quad 32
varSize:
.quad 128
typeSize:
.quad 32 
valSize:
.quad 64 
labelSize:
.quad 128 
labelsMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
buf:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenBuf = . - buf 
buf2:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenBuf2 = . - buf2 
buf3:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenBuf3 = . - buf3
buf4:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenBuf4 = . - buf4
userMem:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserMem = . - userMem 
userMem2:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserMem2 = . - userMem2 
mem:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem = . - mem 
mem2:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem2 = . - mem2 
mem3:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem3 = . - mem3 
mem4:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem4 = . - mem4 
mem5:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem5 = . - mem5
mem6:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem6 = . - mem6
mem7:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem7 = . - mem7
mem8:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem8 = . - mem8
mem9:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem9 = . - mem9
mem10:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem10 = . - mem10 
mem11:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem11 = . - mem11 
mem12:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem12 = . - mem12 
mem13:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem13 = . - mem13 
mem14:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem14 = . - mem14 
mem15:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem15 = . - mem15 
mem16:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem16 = . - mem16 
mem17:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem17 = . - mem17 
mem18:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem18 = . - mem18 
mem19:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMem19 = . - mem19 
strBegin:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrBegin = . - strBegin
oldHeapMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenOldHeapMax = . - oldHeapMax
strPointer:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrPointer = . - strPointer
strMax:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrMax = . - strMax 
isNeg:
.byte 0 
lenIsNeg = . - isNeg
isExpNeg:
.byte 0
lenIsExpNeg = . - isExpNeg
varType:
.quad 0, 0, 0, 0
lenVarType = . - varType 
varName:
.quad 0, 0, 0, 0
lenVarName = . - varName 
userData:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenUserData = . - userData
lastUserData:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenLastUserData = . - lastUserData
strPageNumber:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenStrNumber = . - strPageNumber
memorySize:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMemorySize = . - memorySize
memoryBegin:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenMemoryBegin = . - memoryBegin
labelsEnd:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenLabelsEnd = . - labelsEnd
labelsPointer:
.quad 0, 0, 0, 0, 0, 0, 0, 0
lenLabelsPointer = . - labelsPointer 
systemVarName:
.ascii "^systemVar"
.space 1, 0
lenSystemVarName = . - systemVarName
varName0:
.ascii "iVar"
.space 1, 0
lenVarName0 = . - varName0
varName1:
.ascii "sVar"
.space 1, 0
lenVarName1 = . - varName1
varName2:
.ascii "fVar"
.space 1, 0
lenVarName2 = . - varName2
varName3:
.ascii "bVar"
.space 1, 0
lenVarName3 = . - varName3
varName4:
.ascii "iVar2"
.space 1, 0
lenVarName4 = . - varName4 
varName5:
.ascii "fVar2"
.space 1, 0
lenVarName5 = . - varName5
varName6:
.ascii "sVar2"
.space 1, 0
lenVarName6 = . - varName6
varName7:
.ascii "sVar3"
.space 1, 0
lenVarName7 = . - varName7
varName8:
.ascii "bVar"
.space 1, 0
lenVarName8 = . - varName8
varName9:
.ascii "sVar4"
.space 1, 0
lenVarName9 = . - varName9 

intType:
.ascii "int"
.space 1, 0
lenIntType = . - intType
floatType:
.ascii "float"
.space 1, 0
lenFloatType = . - floatType
boolType:
.ascii "bool"
.space 1, 0
lenBoolType = . - boolType
stringType:
.ascii "string"
.space 1, 0
lenStringType = . - stringType
enter:
.ascii "\n"
.space 1, 0
lenEnter = . - enter 
data0:
.ascii "0"
.space 1, 0
lenData0 = . - data0 
data1:
.ascii "Hello, world!"
.space 1, 0
lenData1 = . - data1 
data2:
.ascii "Slava"
.space 1, 0
lenData2 = . - data2 
data3:
.ascii "3776"
.space 1, 0
lenData3 = . - data3 
data4:
.ascii "5.44"
.space 1, 0
lenData4 = . - data4 
data5:
.ascii "5.44"
.space 1, 0 
lenData5 = . - data5 
data6:
.ascii "1.2"
.space 1, 0 
lenData6 = . - data6 
data7:
.ascii "Slava"
.space 1, 0 
lenData7 = . - data7
data8:
.ascii "CCrypto AG — швейцарская компания, существовавшая в 1952—2018 годах и специализировавшаяся на технологиях в сфере коммуникаций и информационной безопасности. Компанию основал в 1952 году шведский бизнесмен Борис Хагелин. Штаб компании располагался в швейцарском городе Штайнхаузен, в кантоне Цуг. Компания была широко известна в мире как производитель высококачественных шифровальных машин и большого ассортимента различных шифровальных устройств, оставаясь коммерчески успешной на протяжении многих лет. У неё также были офисы в городах Абиджан, Абу-Даби, Буэнос-Айрес, Куала-Лумпур, Маскат и юго-восточном районе Лондона под названием Селсдон. По некоторым данным, у компании были заказчики из 130 разных стран мира, штат компании насчитывал до 400 человек. \n С момента своего основания Crypto AG негласно сотрудничала с американским Агентством национальной безопасности в рамках джентльменского соглашения. В 1970 году Crypto AG в условиях секретности приобрели за 5,75 млн. долларов США две разведслужбы — американская ЦРУ и немецкая БНД, которые владели компанией с 1970 по 1993 год, а позже ЦРУ выкупила немецкую долю и оставалась владельцем Crypto AG вплоть до прекращения её существования в 2018 году. Сведения о владельце фирмы оставались в тайне даже от её менеджеров: право собственности подтверждалось акциями на предъявителя. Обвинения в связях с АНБ и иными западными спецслужбами компания отрицала, несмотря на различные скандалы. \n 11 февраля 2020 года американская газета The Washington Post, немецкая телерадиокомпания ZDF и швейцарская телерадиокомпания Schweizer Radio und Fernsehen официально заявили, что владельцами Crypto AG были именно ЦРУ и БНД, которые с лёгкостью взламывали шифры, сопровождавшие продукты Crypto AG, и свободно расшифровывали закодированные сообщения с 1970 по 2018 год в рамках операции «Рубикон». По итогам швейцарского парламентского расследования выяснилось, что спецслужбы Швейцарии были в курсе оказываемой Crypto AG помощи американским разведслужбам и извлекали из этого выгоду. Публикация информации о том, что ресурсы Crypto AG использовались для взлома зашифрованной переписки высокопоставленных государственных лиц разных стран мира, привела к скандалу. Компанию также обвиняли в продаже программных продуктов с бэкдорами, позволявшими американским, британским и немецким службам радиоэлектронной разведки перехватывать секретную информацию. C"
.space 1, 0
lenData8 = . - data8 
data9:
.ascii "DСтасáнор (IV веке до н. э.) — выходец из аристократической семьи города Солы на Кипре, военачальник Александра Македонского, впоследствии сатрап нескольких провинций на восточной окраине Македонской империи.\nУчаствовал в походах Александра Македонского. Вначале состоял в числе гетайров царя, а затем, после выполнения ответственного поручения Александра, стал сатрапом Арии и Дрангианы. Во время Вавилонского раздела Македонской империи 323 года до н. э. после смерти Александра за Стасанором была подтверждена власть над пожалованными ранее сатрапиями. Во время следующего раздела в Трипарадисе 321/320 года до н. э. он был переведен в Бактрию и Согдиану[fr], которые располагались на северо-восточных окраинах империи. Формальным сатрапом Арии и Дрангианы стал ставленник и, возможно, родственник Стасанора Стасандр. \n После смерти Александра войска Стасанора, как минимум, участвовали в двух военных кампаниях — в составе коалиции сатрапов против Пифона и на стороне Эвмена против Антигона Одноглазого. После поражения Эвмена Стасанор сохранил свою власть, так как новый военный поход вглубь Азии, где располагались Бактрия и Согдиана, не входил в планы Антигона. Однако Стасанор был вынужден формально признать над собой верховную власть Антигона.\nОт подтверждения за Стасанором должности сатрапа Бактрии до завоевания провинции Селевком прошло около десяти лет. Согласно данным нумизматики в этот промежуток времени Стасанора сменил Софит, о биографии которого ничего неизвестно. D"
.space 1, 0
data10:
.ascii "1"
.space 1, 0 
/*data11:
.ascii "CCrypto AG — швейцарская компания, существовавшая в 1952—2018 годах и специализировавшаяся на технологиях в сфере коммуникаций и информационной безопасности. Компанию основал в 1952 году шведский бизнесмен Борис Хагелин. Штаб компании располагался в швейцарском городе Штайнхаузен, в кантоне Цуг. Компания была широко известна в мире как производитель высококачественных шифровальных машин и большого ассортимента различных шифровальных устройств, оставаясь коммерчески успешной на протяжении многих лет. У неё также были офисы в городах Абиджан, Абу-Даби, Буэнос-Айрес, Куала-Лумпур, Маскат и юго-восточном районе Лондона под названием Селсдон. По некоторым данным, у компании были заказчики из 130 разных стран мира, штат компании насчитывал до 400 человек. \n С момента своего основания Crypto AG негласно сотрудничала с американским Агентством национальной безопасности в рамках джентльменского соглашения. В 1970 году Crypto AG в условиях секретности приобрели за 5,75 млн. долларов США две разведслужбы — американская ЦРУ и немецкая БНД, которые владели компанией с 1970 по 1993 год, а позже ЦРУ выкупила немецкую долю и оставалась владельцем Crypto AG вплоть до прекращения её существования в 2018 году. Сведения о владельце фирмы оставались в тайне даже от её менеджеров: право собственности подтверждалось акциями на предъявителя. Обвинения в связях с АНБ и иными западными спецслужбами компания отрицала, несмотря на различные скандалы. \n 11 февраля 2020 года американская газета The Washington Post, немецкая телерадиокомпания ZDF и швейцарская телерадиокомпания Schweizer Radio und Fernsehen официально заявили, что владельцами Crypto AG были именно ЦРУ и БНД, которые с лёгкостью взламывали шифры, сопровождавшие продукты Crypto AG, и свободно расшифровывали закодированные сообщения с 1970 по 2018 год в рамках операции «Рубикон». По итогам швейцарского парламентского расследования выяснилось, что спецслужбы Швейцарии были в курсе оказываемой Crypto AG помощи американским разведслужбам и извлекали из этого выгоду. Публикация информации о том, что ресурсы Crypto AG использовались для взлома зашифрованной переписки высокопоставленных государственных лиц разных стран мира, привела к скандалу. Компанию также обвиняли в продаже программных продуктов с бэкдорами, позволявшими американским, британским и немецким службам радиоэлектронной разведки перехватывать секретную информацию. C"
.space 1, 0 
data12:
.ascii "E«Та сторона, где ветер» — повесть Владислава Крапивина. Состоит из двух частей, которые первоначально были изданы отдельно. Первая часть, «Та сторона, где ветер», была напечатана в журнале «Пионер» в 1965 году, вторая часть, «Люди с фрегата „Африка“» вышла там же в 1967 году. Обе части были объединены в книжном издании 1968 года (издательство «Детская литература») под названием «Та сторона, где ветер» и впоследствии издавались вместе. \n Повесть написана в жанре реализма с элементами романтики. Действие происходит в середине 1960-х годов в городе, который стоит на берегу большой реки. Автор описывает повседневную жизнь нескольких детей, живущих в неполных семьях, их заботы и увлечения, тревоги и мечты. В произведении затронуты темы детской незащищённости, взросления, мужества и самопожертвования, проблемы социальной адаптации ребёнка-инвалида и его возвращения к полноценной жизни. \n Журнальная публикация повести принесла Крапивину известность в СССР. По оценке литературоведа И. Сергиенко, в первой части повести впервые в классической форме возникает образ героя-протагониста, впоследствии получивший название «крапивинские мальчики». Повесть неоднократно переиздавалась и печаталась за рубежом, а в 1978 году была экранизированаE"
.space 1, 0*/
data11:
.ascii "some "
.space 1, 0 
data12:
.ascii "test"
.space 1, 0
data13:
.ascii "FНиканор (др.-греч. Νικάνωρ; IV век до н. э. — 318 или 317 год до н. э.) — военачальник Кассандра, участник Второй войны диадохов. \n В историографии существует несколько взаимоисключающих гипотез, которые идентифицируют Никанора как племянника Аристотеля, либо как сына личного телохранителя-соматофилака македонских царей и сатрапа Киликии[fr] Балакра, либо как одного из военачальников Александра Македонского, который в 334 году до н. э. обеспечил морскую блокаду Милета. \n После смерти регента Македонской империи Антипатра в 319 году до н. э. Никанор принял сторону Кассандра, который боролся за власть в Македонии с Полиперхоном. Никанор руководил македонским гарнизоном в крепости над Афинами Мунихий. Он не только удержал крепость, но и захватил главный порт Аттики Пирей, что стало причиной смены власти в Афинах. Вскоре Кассандр поручил Никанору командование флотом. В сражении при Византии союзный флот Никанора и Антигона Одноглазого разгромил морские силы Полиперхона под командованием Клита Белого. \n Вскоре после победы при Византии Кассандр заподозрил Никанора в измене, приказал арестовать, а затем инициировал войсковое собрание, которое приговорило военачальника к смертиF"
.space 1, 0
data14:
.ascii "переменная только слева"
.space 1, 0
data15:
.ascii "переменная только справа"
.space 1, 0
data16:
.ascii "переменная и слева, и справа"
.space 1, 0 
data17:
.ascii "нет переменных"
.space 1, 0

labelName1:
.ascii "__stop"
.space 1, 0
labelName2:
.ascii "__throughMessage"
.space 1, 0 
ten:
.float 10.0 
one:
.float 1.0 
zero:
.float 0.0 
floatTail:
.ascii ".0"
.space 1, 0 
lenFloatTail =  . - floatTail
label1:
.quad __stop 
label2:
.quad __throughMessage 

fatalError:
.ascii "fatal error: internal error\n"
.space 1, 0 
divZeroError:
.ascii "runtime error: dividing by zero\n"
.space 1, 0
divINegError:
.ascii "runtime error: @ is not defined for negative numbers\n"
.space 1, 0
powNegError:
.ascii "runtime error: ^ is not defined for negative base and fractional exponent\n"
.space 1, 0
powZeroNegError:
.ascii "runtime error: ^ is not defined for zero base and negative exponent\n"
.space 1, 0
powZeroZeroError:
.ascii "runtime error: ^ is not defined for zero base and zero exponent\n"
.space 1, 0
noSuchMarkError:
.ascii "runtime error: no such mark: "
.space 1, 0
concError:
.ascii "could not concatinate not string arguments\n"
.space 1, 0 
strError:
.ascii "the type of the variable to which you want to assign the result of string concatenation is not a string\n"
.space 1, 0 
.text	

__throughError:
 mov $fatalError, %rsi
 call __print 
 mov $60, %rax
 mov $1, %rdi
 syscall

 __throughUserError:
 # %rsi - адрес, по которому лежит сообщение об ошибке 
 call __print 
 mov $60, %rax
 mov $1, %rdi
 syscall

__print:
 mov (%rsi), %al	
 cmp $0, %al	
 jz  __printEx			
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall		    
 inc %rsi		  		    
 jnz __print
__printEx:
 ret

# посчитать количество символов до 0
# адрес начала в %rsi
# результат в %rax  
 __len:
 xor %rax, %rax 
 __lenLocal:
 mov (%rsi), %dl	
 cmp $0, %dl	
 jz  __lenEx				    
 inc %rsi		  	
 inc %rax 	    
 jmp __lenLocal  
__lenEx:
 ret

 __printSymbol:
 mov (%rsi), %al				
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall		    
 ret

__printHeap:  
 mov (memoryBegin), %r8  
 __printHeapLoop:
 cmp (strMax), %r8 
 jz __printHeapEx
 mov (%r8), %dl 
 /*cmp $0, %dl 
 jnz __printHeapNotZero
 mov $endSymbol, %rsi 
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall
 inc %r8 
 jmp __printHeapLoop
 __printHeapNotZero: */
 cmp $2, %dl 
 jnz __printHeapNotTwo
 mov $endSymbol, %rsi 
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall
 inc %r8 
 jmp __printHeapLoop
 __printHeapNotTwo: 

 cmp $1, %dl 
 jnz printHeapNopEnd
 mov $starSymbol, %rsi 
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall
 inc %r8 
 jmp __printHeapLoop
 printHeapNopEnd:
 mov %r8, %rsi 
 mov $1, %rdi	
 mov $1, %rdx
 mov $1, %rax	
 syscall
 inc %r8 
 jmp __printHeapLoop
 __printHeapEx:
 ret 


__toStr:
 # число в %rax 
 # подготовка преобразования числа в строку
  cmp $0, %rax 
  jg __toStrPos
  mov $0, %rdx 
  sub $1, %rdx 
  imul %rdx  
  movb $'-', (buf2)
  jmp __toStrNeg 
  __toStrPos:
  movq $0, (buf2)
  __toStrNeg:
  mov $10, %r8    # делитель
  mov $buf, %rsi  # адрес начала буфера 
  xor %rdi, %rdi  # обнуляем счетчик
# преобразуем путем последовательного деления на 10
__toStrlo:
  xor %rdx, %rdx  # число в rdx:rax
  div %r8         # делим rdx:rax на r8
  add $48, %dl    # цифру в символ цифры
  mov %dl, (%rsi) # в буфер
  inc %rsi        # на следующую позицию в буфере
  inc %rdi        # счетчик увеличиваем на 1
  cmp $0, %rax    # проверка конца алгоритма
  jnz __toStrlo          # продолжим цикл?
# число записано в обратном порядке,
# вернем правильный, перенеся в другой буфер 
  mov $buf2, %rbx # начало нового буфера
  mov (%rbx), %dl 
  cmp $'-', %dl
  jnz __toStrEmpty
  inc %rbx 
  __toStrEmpty:
  mov $buf, %rcx  # старый буфер
  add %rdi, %rcx  # в конец
  dec %rcx        # старого буфера
  mov %rdi, %rdx  # длина буфера
# перенос из одного буфера в другой
__toStrexc:
  mov (%rcx), %al # из старого буфера
  mov %al, (%rbx) # в новый
  dec %rcx        # в обратном порядке  
  inc %rbx        # продвигаемся в новом буфере
  dec %rdx        # а в старом в обратном порядке
  jnz __toStrexc         # проверка конца алгоритма 
  movb $0, (%rbx)
  mov $buf2, %rbx 
  mov (%rbx), %dl 
  cmp $'-', %dl
  jnz __toStrEnd 
  inc %rbx 
  mov (%rbx), %dl 
  cmp $'0', %dl
  jnz __toStrEnd 
  inc %rbx 
  mov (%rbx), %dl 
  cmp $0, %dl 
  jnz __toStrEnd
  mov $buf2, %rbx 
  movb $'0', (%rbx) 
  inc %rbx 
  movb $0, (%rbx) 
__toStrEnd:
  ret
  
__set: #set strings
 # входные параметры 
 # rsi - длина буфера назначения 
 # rdx - адрес буфера назначения
 # rax - длина буфера источника 
 # rdi - адрес буфера источника 

 mov %rdx, %r8 
 mov %rsi, %r9
  
 __setClear:
 cmp $0, %rsi
 jz __setClearEnd
 movb $1, (%rdx)
 dec %rsi
 inc %rdx  
 jmp __setClear
 __setClearEnd:
 dec %rdx 
 movb $0, (%rdx)

 mov %r8, %rdx  
 mov %r9, %rsi  
 
 __setLocal:
 cmp $0, %rax 
 jz __setLocalEnd
 cmp $2, %rsi
 jz __setLocalEnd
 mov (%rdi), %r11b 
 movb %r11b, (%rdx)
 inc %rdx
 inc %rdi
 dec %rax  
 dec %rsi 
 jmp __setLocal
 __setLocalEnd:
 dec %rdx 
 mov (%rdx), %rax 
 cmp $1, %rax 
 jz __star
 inc %rdx
 __star: 
 movb $0, (%rdx)
 ret 

__concatinate:
 # входные параметры 
 # r8 - длина буфера первого операнда 
 # r9 - адрес буфера первого операнда
 # r11 - адрес буфера второго операнда 
 # выход
 # userData   
 call __clearUserData
 mov $lenUserData, %rsi # присваиваем в userData первый операнд  
 mov $userData, %rdx 
 mov %r8, %rax 
 mov %r9, %rdi 

 mov %r8, %rbx 
 mov %r11, %rcx
 call __set
 mov %rbx, %r8 
 mov %rcx, %r11 

 mov $userData, %r8
 __concNext:
 mov (%r8), %dl 
 cmp $0, %dl 
 jz __concLocal 
 inc %r8   
 jmp __concNext 
  
 __concLocal:
 
 mov (%r11), %dl 
 movb %dl, (%r8)
 inc %r8 
 inc %r11
 #dec %r10  
 cmp $0, %dl 
 jnz __concLocal
 #movb $0, (%r8)
 
 ret   
 
 __userConcatinateVars:
 # входные параметры 
 # r8 - адрес имени первой переменной 
 # r9 - адрес имени второй переменной
 # $varName - адрес имени переменной, куда положить результат 

 #mem13 - 16 свободны 
 mov $varName, %rsi 
 call __print  

 ret 

 __userConcatinate:
 # входные параметры 
 # r8 - адрес начала первой строки (либо адрес имени переменной)
 # r9 - адрес начала второй строки (либо адрес имени переменной)
 # rax = 1 - строка слева лежит в переменной
 # rax = 0 - строка слева в статической памяти
 # rbx = 1 - строка справа лежит в переменной 
 # rbx = 0 - строка справа лежит в статичесой памяти 
 # $varName - адрес имени переменной, куда положить результат
 mov %r8, (mem13)
 mov %r9, (mem14)
 mov %rax, (mem15)
 mov %rbx, (mem16)
 
 # сохранили приемник
 mov $lenMem17, %rsi 
 mov $mem17, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set

 mov $mem17, %r8  
 mov %r8, (userData)
 // сохраним в системной переменной значение, которое сейчас хранится в приемнике
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName, %rax 
 mov $systemVarName, %rdi 
 call __set

 mov $1, %rax 
 call __setVar 

 # восстановили приемник
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenMem17, %rax 
 mov $mem17, %rdi 
 call __set

 call __getVar

 mov %r13, %rbx
 __userConcatinateLocal:
 cmp %r15, %rbx
 jg __userConcatinateEnd

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __userConcatinateEnd 
 
 add (varSize), %rbx 
 jmp __userConcatinateLocal  
  
 __userConcatinateEnd:
 __userConcatinateSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __throughError
 mov $buf, %rsi 
 mov %rbx, %r12 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jz __userConcatinateSearch
 

 add (varNameSize), %rbx 
 mov %rbx, %rax 
 mov %rbx, %r12 
 call __read  
 add (typeSize), %rbx 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set
 mov %rbx, %r12 
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jnz __userConcatinateErrEnd
 mov $strError, %rsi 
 call __throughUserError

 __userConcatinateErrEnd:

 
 mov (userData), %rbx 
 mov (mem13), %rax 
 
 __userConcatinateClearStr:
 mov (%rbx), %dil 
 cmp $2, %dil 
 jz __userConcatinateClearStrEnd
 movb $1, (%rbx)
 inc %rbx 
 jmp __userConcatinateClearStr
 __userConcatinateClearStrEnd:
 dec %rbx 
 movb $0, (%rbx)
 mov (userData), %rbx 

 __userConcatinateIsNotStr:
 
 cmp $1, (mem15)
 jnz __userConcatinateRightVar
 cmp $1, (mem16)
 jz __userConcatinateTwoVars
 
 // переменная только слева 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenMem13, %rax 
 mov (mem13), %rdi 
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare 

 cmp $1, %rax 
 jnz __userConcatinateNotTheSameLeft
 
 mov $systemVarName, %rax 
 mov %rax, (userData)

 mov $1, %rax 

 call __setVar 
 jmp __userConcatinateTheSameLeft

 __userConcatinateNotTheSameLeft:

 mov (mem13), %rsi 
 mov %rsi, (userData)
 mov $1, %rax 
 
 call __setVar
 
 __userConcatinateTheSameLeft:
 call __getVar 
 mov (userData), %rsi 
 call __len 
 mov %rax, %rbx 
 add (userData), %rbx 
 mov (mem14), %rax 

 jmp __userConcatinateNow2 

 __userConcatinateTwoVars:
 // с обеих сторон переменная
 mov (mem13), %r8 
 mov (mem14), %r9 
 call __userConcatinateVars 
 ret 

 __userConcatinateRightVar:
 cmp $1, (mem16)
 jnz __userConcatinateNoVars
 // переменная только справа 
 mov (mem13), %rsi 
 call __len 
 mov %rax, (mem16) # сохранили длину первого операнда
 
 # сохранили приемник
 mov $lenMem15, %rsi 
 mov $mem15, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov (mem14), %rdi 
 call __set
 call __getVar 

 # восстановили приемник
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenMem15, %rax 
 mov $mem15, %rdi 
 call __set
 
 mov (userData), %rsi 
 call __len 
 add (mem16), %rax # длина результата 
 mov %rax, (mem16)

 call __getVar 
 mov (userData), %rbx
 mov (mem16), %rax  


  __userConcatinatePrepare:
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __userConcatinateMoreMemEnd0
 mov %rax, (mem4)
 mov %rbx, (mem5) 
 mov %r12, (mem10)
 call __internalShiftStr
 mov (mem10), %r12 
 mov (mem4), %rax
 mov (mem5), %rbx
 __userConcatinateMoreMemEnd0:
 cmp $0, %rax 
 jz __userConcatinatePrepareEnd
 inc %rbx 
 dec %rax 
 jmp __userConcatinatePrepare
 __userConcatinatePrepareEnd:

 call __getVar 
 mov (userData), %rbx

 mov (mem13), %rax

 __userConcatinateNow0: 
 mov (%rax), %dl
 cmp $0, %dl 
 jz __userConcatinateRet0 
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 
 jmp __userConcatinateNow0 
 __userConcatinateRet0:

 mov %rbx, (mem16) # сохранили %rbx, куда нужно записывать результат

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName, %rax 
 mov (mem14), %rdi 
 call __set
 call __getVar 

 mov (userData), %rax 
 mov (mem16), %rbx 

 __userConcatinateNow1: 
 mov (%rax), %dl
 cmp $0, %dl 
 jz __userConcatinateRet1 
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 
 jmp __userConcatinateNow1 
 __userConcatinateRet1:

 ret  

 __userConcatinateNoVars:
 // нет переменных 
 
 __userConcatinateEndCheck:
 __userConcatinateNow:
  
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __userConcatinateMoreMemEnd

 mov %rax, (mem4)
 mov %rbx, (mem5) 
 mov %r12, (mem10)
 call __internalShiftStr
 mov (mem10), %r12 
 mov (mem4), %rax
 mov (mem5), %rbx 
 
 __userConcatinateMoreMemEnd:
 mov (%rax), %dl
 cmp $0, %dl 
 jz __userConcatinateRet 
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 
 jmp __userConcatinateNow 

 __userConcatinateRet: 
 
 mov (mem14), %rax  
__userConcatinateNow2:
   
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __userConcatinateMoreMemEnd2

 mov %rax, (mem4)
 mov %rbx, (mem5) 
 mov %r12, (mem10)
 call __internalShiftStr
 mov (mem10), %r12 
 mov (mem4), %rax
 mov (mem5), %rbx 
 
 __userConcatinateMoreMemEnd2:
 mov (%rax), %dl
 cmp $0, %dl 
 jz __userConcatinateRet2 
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 
 jmp __userConcatinateNow2 

 __userConcatinateRet2:
 movb $0, (%rbx) 

 mov (mem13), %r8 
 mov (mem14), %r9 
 mov (mem15), %rax 
 mov (mem16), %rbx 
 ret
 

__toNumber:
 # вход: buf 
 # выход:  %rax 
 mov $buf, %rdx # our string
 movzx (%rdx), %rcx 
 cmp $'-', %rcx 
 jnz __toNumberAtoi
 inc %rdx 
 __toNumberAtoi:
 xor %rax, %rax # zero a "result so far"
 __toNumberTop:
 movzx (%rdx), %rcx # get a character
 inc %rdx # ready for next one
 cmp $0, %rcx # end?
 jz __toNumberDone
 sub $48, %rcx # "convert" character to number
 imul $10, %rax # multiply "result so far" by ten
 add %rcx, %rax # add in current digit
 jmp __toNumberTop # until done
 __toNumberDone:
 mov $buf, %rdx 
 movzx (%rdx), %rcx 
 cmp $'-', %rcx 
 jnz __toNumberIsPos
 mov $0, %rdx 
 sub $1, %rdx 
 //mov %rdx, (buf)
 imul %rdx 
 __toNumberIsPos:
 ret


__defineVar:
 # адрес имени переменной в $varName
 # адрес типа переменой в $varType
 
 mov $varName, %rcx 
 mov $varType, %rdx 
 cmp %r14, %r15
 jg __defOk 
 #mov %r15, %r8
 mov (strMax), %r8 
 call __newMem  
 mov $varName, %rcx 
 mov $varType, %rdx
 __defOk:
 mov %r14, %r8 
 __defOkLocal: 
 movb (%rcx), %r11b 
 cmp $1, %r11b
 jz __defOkLocalEx
 movb %r11b, (%r8)
 inc %rcx 
 inc %r8 
 jmp __defOkLocal
 __defOkLocalEx:
 mov %r14, %r8 
 add (varNameSize), %r8 
 __defOkTypeLocal:
 movb (%rdx), %r11b
 cmp $1, %r11b 
 jz __defOkTypeLocalEx
 movb %r11b, (%r8)
 inc %rdx
 inc %r8 
 jmp __defOkTypeLocal
 __defOkTypeLocalEx:
 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarType, %rax 
 mov $varType, %rdi 
 call __set 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi 
 call __set 
  
 call __compare 
 cmp $1, %rax 
 jz __defInt 

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenFloatType, %rax 
 mov $floatType, %rdi 
 call __set 
 call __compare
 cmp $1, %rax 
 jz __defFloat 

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi 
 call __set 
 call __compare
 cmp $1, %rax  
 jz __defBool 

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set 
 call __compare
 cmp $1, %rax 
 jz __defString 
 call __throughError

 __defInt:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 movb $'0',(%r8)
 inc %r8 
 movb $0, (%r8)
 jmp __defEnd
 __defFloat:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 movb $'0',(%r8)
 inc %r8 
 movb $'.', (%r8)
 inc %r8 
 movb $'0', (%r8)
 inc %r8 
 movb $0, (%r8)
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8 
 jmp __defEnd
 __defBool:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 movb $'0',(%r8)
 inc %r8 
 movb $0, (%r8) 
 jmp __defEnd
 __defString:
 mov %r14, %r8 
 add (varNameSize), %r8 
 add (typeSize), %r8
 mov (strPointer), %rax 
 movb $0, (%rax)
 add (valSize), %rax
 dec %rax  
 movb $2, (%rax) # признак конца поля для строки 
 inc %rax 
 sub (valSize), %rax 
 mov %r8, %r12 
 call __toStr 
 mov %r12, %r8 
 mov $buf2, %rax
 __defAddr:
 mov (%rax), %dl 
 cmp $0, %dl  
 jz __defStrEnd 
 mov %dl, (%r8)
 inc %rax 
 inc %r8 
 jmp __defAddr
 __defStrEnd:
 movb $0, (%r8)
 mov (strPointer), %rax 
 add (valSize), %rax 
 #cmp (strMax), %rax 
 #jg __defStrNewMem
 mov %rax, (strPointer)
 #jmp __defEnd 
 #__defStrNewMem:
 #mov %rax, (strPointer)
 #call __newStrMem
 __defEnd:

 add (varSize), %r14
 ret 

__undefineVar:
 # вход: имя переменной по адресу $varName 
 mov %r13, %rbx
 __undefVarLocal:
 cmp %r15, %rbx
 jg __undefVarEnd

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __undefVarEnd 
 
 add (varSize), %rbx 
 jmp __undefVarLocal  
  
 __undefVarEnd:
 __undefVarSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __undefEnd
 mov $buf, %rsi 
 mov %rbx, %r12 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jz __undefVarSearch
 
 mov %rbx, %r12 
 add (varNameSize), %r12
 __undefName: 
 cmp %rbx, %r12 
 jz __undefNameEx
 movb $'!', (%rbx)
 inc %rbx 
 jmp __undefName 
 __undefNameEx:
 dec %rbx 
 movb $0, (%rbx)
 inc %rbx 
 mov %rbx, %r12 
 add (typeSize), %r12 
 __undefType: 
 cmp %rbx, %r12 
 jz __undefTypeEx
 movb $'!', (%rbx)
 inc %rbx 
 jmp __undefType 
 __undefTypeEx:
 dec %rbx
 dec %rbx  
 movb $0, (%rbx)
 inc %rbx 
 movb $0, (%rbx)
 inc %rbx 
 mov %rbx, %r12 
 add (valSize), %r12 
 __undefVal: 
 cmp %rbx, %r12 
 jz __undefValEx
 movb $'!', (%rbx)
 inc %rbx 
 jmp __undefVal  
 __undefValEx:
 dec %rbx 
 movb $0, (%rbx)
 __undefEnd:
 ret 

# r12 - pointer (общего назначения)
# r13 - heapBegin 
# r14 - heapPointer 
# r15 - heapMax 

__firstMem:
 # в %rax адрес начала выделяемой памяти 
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8  
 mov %rax, %r13
 mov %rax, %r14
 mov %rax, %r9 
 add (pageSize), %r9
 mov %r9, %r15 
# выделить динамическую память
 mov (pageSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $1, %dl
 mov $0, %rbx
 __lo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (pageSize), %rbx
 jz  __ex
 jmp __lo
 __ex:
 ret 

 __newLabelMem:
 # в %rax адрес начала выделяемой памяти 
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8  
 mov %rax, %r9 
 add (labelSize), %r9
# выделить динамическую память
 mov (labelSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $1, %dl
 mov $0, %rbx
 __newLabelMemlo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (labelSize), %rbx
 jz  __newLabelex
 jmp __newLabelMemlo
 __newLabelex:
 ret

 __firstStrMem:
 # адрес начала области для выделения памяти
 mov %r15, %rax
# запомнить адрес начала выделяемой памяти
 mov %rax, %r8  
 mov %rax, (strBegin)
 mov %rax, (strPointer)
 mov %rax, %r9 
 add (pageSize), %r9
 mov %r9, (strMax)
# выделить динамическую память
 mov (pageSize), %rdi
 add %rax, %rdi
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $1, %dl
 mov $0, %rbx
 __firstStrMemLo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (pageSize), %rbx
 jz  __firstStrMemEx
 jmp __firstStrMemLo
 __firstStrMemEx:
 ret 


 __newMem:
 # адрес начала выделяемой памяти в  %r8 
# запомнить адрес начала выделяемой памяти
 #mov %r8, %r14 
 mov %r8, %r9 
 add (pageSize), %r9 
 #mov %r9, %r15
 mov %r15, (oldHeapMax)
 #mov (strPointer), %r15 
 add (pageSize), %r15 
# выделить динамическую память
 mov (pageSize), %rdi
 #add %rax, %rdi
 add %r8, %rdi 
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $1, %dl
 mov $0, %rbx
 __newMemlo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (pageSize), %rbx
 jz  __newMemEx
 jmp __newMemlo
 __newMemEx:
 mov %r15, %r8 
 call __newStrMem
 call __shiftStr
 ret  

 __newStrMem:
 # адрес начала выделяемой памяти в  %r8 
# запомнить адрес начала выделяемой памяти
 #mov %r8, %r14
 mov %r8, %r9 
 #add (pageSize), %r9
 
 mov (strPageNumber), %rax
 mov (pageSize), %rdi
 __newStrMemOkBegin: 
 cmp $0, %rax 
 jz __newStrMemOk
 dec %rax 
 add (pageSize), %rdi  
 jmp __newStrMemOkBegin
 __newStrMemOk:
 mov %rdi, (memorySize)
 add (memorySize), %r9 
 mov (strPageNumber), %rax 
 inc %rax 
 mov %rax, (strPageNumber)
 #mov %r9, (strMax)
# выделить динамическую память
 mov (memorySize), %rdi
 #add %rax, %rdi
 add %r8, %rdi 
 mov $12, %rax
 syscall
# обработка ошибки
 cmp $-1, %rax
 jz __throughError
# заполним выделенную память
 mov $1, %dl
 mov $0, %rbx
 __newStrMemlo:
 movb %dl, (%r8)
 inc %rbx
 inc %r8 
 cmp (memorySize), %rbx
 jz  __newStrMemEx
 jmp __newStrMemlo
 __newStrMemEx:
 ret 

__read: 
 # считать в buf по указателю в %r12 
 mov %r12, %r8
 mov $buf, %r10 
 mov $lenBuf, %rsi 
__readClear:
 movb $0, (%r10)
 dec %rsi
 inc %r10
 cmp $0, %rsi  
 jnz __readClear
 mov $buf, %r10
 __readLocal: 
 movb (%r8), %r9b
 cmp $1, %r9b  
 jz __readEx
 cmp $0, %r9b  
 jz __readEx
 mov %r9b, (%r10)
 inc %r10 
 inc %r8 
 jmp __readLocal
 __readEx:
 mov $buf, %r10 
 cmp $0, (%r10)
 jnz __readOk
 movb $1, (%r10)
 __readOk:
 ret 
 
 __renewStr:
 # адрес начала кучи 
 mov %r13, %r12
 # старый адрес конца кучи 
 mov (oldHeapMax), %r11 
 add (varNameSize), %r12  
 
 __renewFindStr:
 call __read
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set
 call __compare 
 cmp $1, %rax 
 jz __renewVal
 add (varSize), %r12 
 cmp %r11, %r12  
 jge __renewStrEnd
 jmp __renewFindStr

 __renewVal:
 add (typeSize), %r12 
 call __read 
 call __toNumber
 __renewValLocal:
 mov (%rax), %r10b 
 cmp $0, %r10b 
 jz __renewValEnd
 movb $1, (%rax)
 inc %rax 
 jmp __renewValLocal 
 __renewValEnd: 
 movb $1, (%rax)
 __renewAddr:
 call __read 
 mov $buf, %rsi
 call __toNumber 
 #mov (pageSize), %rax 
 #call __toStr
 add (pageSize), %rax 
 call __toStr # в buf2 новый адрес строки 

 mov %r12, %rsi 
 __renewAddrLocal:
 mov (%rsi), %r10b 
 cmp $1, %r10b 
 jz __renewAddrEnd
 movb $1, (%rsi)
 inc %rsi 
 jmp __renewAddrLocal
 __renewAddrEnd:
 mov %r12, %rsi 
 mov $buf2, %rdx 
 // запись нового адреса
 __renewStrAddr: 
 mov (%rdx), %r10b 
 cmp $0, %r10b 
 jz __renewStrAddrEnd
 movb %r10b, (%rsi)
 inc %rsi 
 inc %rdx 
 jmp __renewStrAddr
 __renewStrAddrEnd:
 movb $0, (%rsi)
 sub (typeSize), %r12 
 add (varSize), %r12 
 jmp __renewFindStr
 __renewStrEnd:
  
 ret 

 __shiftStr:
 # формируем в %r10 адрес нового начала
 mov (strBegin), %r10 
 add (pageSize), %r10
 #mov %r10, %r12 
 #add (pageSize), %r12 
 #mov %r12, (strMax)
 # адрес старого начала
 mov (strBegin), %r11
 __shiftMake: 
 mov (strMax), %r9
 cmp %r9, %r11   
 jz __shiftMakeEnd
 movb (%r11), %r12b 
 movb %r12b, (%r10)
 inc %r10
 inc %r11 
 jmp __shiftMake
 __shiftMakeEnd:
 mov (strPointer), %r10 
 add (pageSize), %r10 
 mov %r10, (strPointer)

 mov (strBegin), %r10 
 add (pageSize), %r10 
 mov %r10, (strBegin)

 mov (strMax), %r10 
 add (pageSize), %r10 
 mov %r10, (strMax)
 
 call __renewStr
 ret 

__compare:
 # сравнить строки по адресу $buf и $buf2
 # если длины строк не равны, то строки не равны 
 mov $buf, %rsi 
 call __len 
 mov %rax, %rbx 
 mov $buf2, %rsi 
 call __len 
 cmp %rax, %rbx 
 jnz __notEqual

 mov $buf, %rax 
 mov $buf2, %rbx 
 __compareLocal:
 movb (%rax), %dl 
 cmp $0, %dl 
 jz __equal
 movb (%rax), %dl
 cmp %dl, (%rbx) 
 jnz __notEqual
 inc %rax 
 inc %rbx 
 jmp __compareLocal

 __notEqual:
 mov $0, %rax 
 ret 
 __equal:
 mov $1, %rax  
 ret 

__internalMakeShiftStr:
# %r12 - адрес внутри таблицы строк, начиная с которого нужно сделать сдвиг 
 mov %r12, %rax # начиная с этого адреса нужно сдвинуть непосредственно строки на valSize   
 mov %rax, (mem11)
 # формируем адрес нового конца 
 mov (strPointer), %r10 
 add (valSize), %r10
 mov %r10, (mem12)
 cmp (strMax), %r10 
 jl __internalMakeShiftStrOk
 mov (strMax), %r8 
 call __newStrMem
 mov (pageSize), %r8 
 add %r8, (strMax) 
 __internalMakeShiftStrOk:
 # адрес старого конца
 mov (strPointer), %r11
 mov (mem12), %r10 
 mov (mem11), %rax  
 mov %rax, %r9
 __internalShiftMake: 

 cmp %r9, %r11   
 jl __internalShiftMakeEnd
 movb (%r11), %al 
 movb %al, (%r10)
 movb $1, (%r11)

 dec %r10
 dec %r11 
 jmp __internalShiftMake
 __internalShiftMakeEnd:
 
 ret 


__internalShiftStr:
# %r12 - место внутри таблицы переменных, после которого нужно сделать сдвиг  
mov (valSize), %rsi 
add %rsi, (strPointer) 

mov %r12, (mem9)
call __read

call __toNumber
mov (mem9), %r12 

mov %rax, (mem9) # с этого адреса нужно сделать сдвиг в таблице строк 

mov %r12, %rsi 
sub (typeSize), %rsi 
add (varSize), %rsi

__internalShiftStrLocal: 
cmp %r14, %rsi  
 
jge __internalShiftStrEnd

mov %rsi, (mem)
mov %rax, (mem2)

mov %rsi, %rdi  
call __len 
mov $lenBuf2, %rsi 
mov $buf2, %rdx 
mov $lenVarType, %rax 
call __set 

mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenStringType, %rax 
mov $stringType, %rdi 
call __set 

call __compare 
cmp $1, %rax 

jnz __internalShiftStrChangeNext

mov (mem), %rsi 
mov (mem2), %rax

mov (mem), %rdi 
add (typeSize), %rsi 
add (typeSize), %rdi 
add (valSize), %rdi 
mov %rdi, (mem6)
mov %rsi, (mem7)
mov %rsi, %r12 
call __read 
call __toNumber
mov %rax, (mem8)
add (valSize), %rax 

call __toStr 
mov (mem6), %rdi
mov (mem7), %rsi   
__internalShiftStrClear:
cmp %rdi, %rsi 
jge __internalShiftStrClearEnd
movb $1, (%rsi)
inc %rsi 
jmp __internalShiftStrClear

__internalShiftStrClearEnd:
dec %rsi 
#movb $0, (%rsi) 

mov (mem), %rsi 
add (typeSize), %rsi 
 
mov $buf2, %rdi 
 
__internalShiftStrSet:
mov (%rdi), %al 
cmp $0, %al 
jz __internalShiftStrChangeEnd
mov (%rdi), %al 
mov %al, (%rsi)
inc %rdi 
inc %rsi 
jmp __internalShiftStrSet

__internalShiftStrChangeEnd:
movb $0, (%rsi)

__internalShiftStrChangeNext:
mov (mem), %rsi 
mov (mem2), %rax
add (varSize), %rsi 

mov %rsi, (mem)
mov %rax, (mem2)

mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenBuf2, %rax 
mov $buf2, %rdi 
call __set 

call __toNumber
cmp (strMax), %rax 
mov (mem), %rsi 
mov (mem2), %rax 
jl __internalShiftStrOk 
mov (strMax), %r8 
call __newStrMem
mov (pageSize), %r8
add %r8, (strMax)
mov (mem), %rsi 
mov (mem2), %rax 
 
__internalShiftStrOk:

jmp  __internalShiftStrLocal

__internalShiftStrEnd:
mov (mem9), %r12 

__internalMakeShiftStrNowLoop:
mov (%r12), %dl 
cmp $2, %dl 
jz __internalMakeShiftStrNow
inc %r12 

jmp __internalMakeShiftStrNowLoop
__internalMakeShiftStrNow:


call __internalMakeShiftStr
ret 


__setVar:
 # вход: 
 # имя переменной по адресу $varName 
 # данные по указателю в (userData) - статика 
 # адрес имени переменной в (userData) - данные в переменной 
 # rax = 0 - статические данные 
 # rax = 1 - данные в переменной
 cmp $1, %rax 
 jnz __setVarStat

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf, %rax 
 mov (userData), %rdi 
 call __set
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare 

 cmp $1, %rax 
 jnz __setVarStat0 
 ret # присвоение переменной в саму себя 
 __setVarStat0:
 mov $1, %rax 
 __setVarStat:
 mov %rax, (mem4) 

 mov %r13, %rbx
 __setVarLocal:
 cmp %r15, %rbx
 jg __setVarEnd

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __setVarEnd 
 
 add (varSize), %rbx 
 jmp __setVarLocal  
  
 __setVarEnd:
 __setVarSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __throughError
 mov $buf, %rsi 
 mov %rbx, %r12 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jz __setVarSearch
 

 add (varNameSize), %rbx 
 mov %rbx, %rax 
 mov %rbx, %r12 
 call __read  
 add (typeSize), %rbx 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set
 mov %rbx, %r12 
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jnz __setVarStr
 
 #add (varNameSize), %rbx 
 #add (typeSize), %rbx
 mov %rbx, %r8 
 mov %rbx, %rax 
 __setVarClear:
 add (valSize), %rax 
 __setVarClearLocal: 
 cmp %rax, %rbx 
 jz __setVarClearEnd
 movb $1, (%rbx) 
 inc %rbx 
 jmp __setVarClearLocal
 __setVarClearEnd:
 mov %r8, %rbx 
 jmp __setVarIsNotStr

 __setVarStr:
 mov %rbx, %r12 
 call __read 
 call __toNumber 
 mov %rax, %rbx  
 mov %rbx, %r10 

 __setVarClearStr:
 mov (%rbx), %dil 
 cmp $2, %dil 
 jz __setVarClearStrEnd
 movb $1, (%rbx)
 inc %rbx 
 jmp __setVarClearStr
 __setVarClearStrEnd:
 dec %rbx 
 movb $0, (%rbx)
 mov %r10, %rbx 
 mov %rbx, (mem5)
 mov %r12, (mem10)

 __setVarIsNotStr:

 cmp $1, (mem4) # переменная? 
 jnz __setVarNotVar
 
 mov %rbx, (mem19) # сохраним %rbx 

 mov $lenMem18, %rsi # сохранить varName 
 mov $mem18, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenUserData, %rax 
 mov (userData), %rdi 
 call __set
 
 call __getVar
 

 mov (userData), %rsi 
 call __len  
 mov (mem5), %rbx 
 mov (mem10), %r12 

 __setVarPrepare:
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __setVarMoreMemEnd0
 mov %rax, (mem4)
 mov %rbx, (mem5) 
 mov %r12, (mem10)
 call __internalShiftStr
 mov (mem10), %r12 
 mov (mem4), %rax
 mov (mem5), %rbx
 __setVarMoreMemEnd0:
 cmp $0, %rax 
 jz __setVarPrepareEnd
 inc %rbx 
 dec %rax 
 jmp __setVarPrepare
 __setVarPrepareEnd:
 call __getVar # обновим указатель в (userData) 
 mov (mem19), %rbx # восстановим первоначальный %rbx 

 __setVarNotStr0:
 mov (userData), %rax 
 __setNow0:
 mov (%rax), %dl
 cmp $0, %dl 
 jz __setVarRet0
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 jmp __setNow0

 __setVarRet0:
 movb $0, (%rbx)
 mov $lenVarName, %rsi # восстановим varName 
 mov $varName, %rdx 
 mov $lenMem18, %rax 
 mov $mem18, %rdi 
 call __set
 ret
 __setVarNotVar: 
 mov (userData), %rax 
 __setNow:
   
 mov (%rbx), %dil 
 cmp $2, %dil 
 jnz __setVarMoreMemEnd
 mov %rax, (mem4)
 mov %rbx, (mem5) 
 mov %r12, (mem10)
 call __internalShiftStr
 mov (mem10), %r12 
 mov (mem4), %rax
 mov (mem5), %rbx 
 
 __setVarMoreMemEnd:
 mov (%rax), %dl
 cmp $0, %dl 
 jz __setVarRet 
 mov %dl, (%rbx)
 inc %rbx 
 inc %rax 
 
 jmp __setNow 

 __setVarRet:
 movb $0, (%rbx) 

 ret

 __getVar:
 # вход: имя переменной по адресу $varName 
 # выход: указатель на данные в (userData)
 mov %r13, %rbx
 __getVarLocal:
 cmp %r15, %rbx
 jg __getVarEnd

 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __getVarEnd 
 
 add (varSize), %rbx 
 jmp __getVarLocal  
  
 __getVarEnd:
 __getVarSearch:
 sub (varSize), %rbx 
 mov %rbx, %r12 
 call __read 
 cmp $1, (buf)
 jz __throughError
 mov $buf, %rsi 
 mov %rbx, %r12 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenVarName, %rax 
 mov $varName, %rdi 
 call __set
 call __compare
 mov %r12, %rbx 
 cmp $0, %rax 
 jz __getVarSearch
 
 add (varNameSize), %rbx 
 mov %rbx, %r12 
 call __read 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi 
 call __set
 call __compare 
 cmp $1, %rax 
 jnz __getVarNotStr
 mov $1, %r11
 jmp __getVarNotStrEnd
 __getVarNotStr:
 mov $0, %r11 
 __getVarNotStrEnd:
 mov %r12, %rbx 
 add (typeSize), %rbx  
 
 __getNow:
 #call __toNumber
 cmp $1, %r11
 jz __getVarGetStr
 mov %rbx, (userData)
 ret 
 __getVarGetStr:
 mov %rbx, %r12 
 call __read
 call __toNumber
 mov %rax, (userData) 
 ret 
 /*mov $userData, %rsi
 __getNowLocal:  
 cmp $0, %rax 
 jz __getVarRet
 mov (%rbx), %dl 
 mov %dl, (%rsi)
 inc %rsi 
 inc %rbx  
 dec %rax 
 jmp __getNowLocal 
 __getVarRet:
 ret*/

__clearBuf:
mov $buf, %rsi 
mov $lenBuf, %rdi
__clearBufLocal: 
cmp $1, %rdi 
jz __clearBufEnd
movb $1, (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal

__clearBufEnd:
movb $0, (%rsi)
ret

__clearBuf2:
mov $buf2, %rsi 
mov $lenBuf2, %rdi
__clearBufLocal2: 
cmp $1, %rdi 
jz __clearBufEnd2
movb $1, (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal2

__clearBufEnd2:
movb $0, (%rsi)
ret

__clearUserData:
mov $userData, %rsi 
mov $lenUserData, %rdi
__clearUserDataLocal: 
cmp $1, %rdi 
jz __clearUserDataEnd
movb $1, (%rsi)
inc %rsi 
dec %rdi 
jmp __clearUserDataLocal

__clearUserDataEnd:
movb $0, (%rsi)
ret

__clearBuf3:
mov $buf3, %rsi 
mov $lenBuf3, %rdi
__clearBufLocal3: 
cmp $1, %rdi 
jz __clearBufEnd3
movb $1, (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal3

__clearBufEnd3:
movb $0, (%rsi)
ret

__clearBuf4:
mov $buf4, %rsi 
mov $lenBuf4, %rdi
__clearBufLocal4: 
cmp $1, %rdi 
jz __clearBufEnd4
movb $1, (%rsi)
inc %rsi 
dec %rdi 
jmp __clearBufLocal4

__clearBufEnd4:
movb $0, (%rsi)
ret

__add:
 # вход: buf и buf2
 # %rax - тип операции 
 # 0 - целочисленное сложение 
 # 1 - сложение вещественных чисел   
 # выход: userData 
 call __clearUserData
 cmp $0, %rax 
 jz __addInt 
 cmp $1, %rax 
 jz __addFloat

 call __throughError

 __addInt:
 call __toNumber
 mov %rax, %rbx 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 call __toNumber
 add %rbx, %rax
 call __toStr 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 ret 
 __addFloat:
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, (buf)
 fld (buf)
 movss %xmm1, (buf)
 fadd (buf)
 fstp (buf)
 call __floatToStr

 ret 


 __sub:
 # вход: buf и buf2
 # %rax - тип операции 
 # 0 - целочисленное сложение 
 # 1 - сложение вещественных чисел   
 # выход: userData 
 call __clearUserData
 cmp $0, %rax 
 jz __subInt 
 cmp $1, %rax 
 jz __subFloat

 call __throughError

 __subInt:
 call __toNumber
 mov %rax, %rbx 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 call __toNumber
 sub %rax, %rbx
 mov %rbx, %rax 
 call __toStr 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 ret 
 __subFloat:
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 movss %xmm1, (buf)
 fld (buf)
 movss %xmm0, (buf)
 fsub (buf)
 fstp (buf)
 call __floatToStr

 ret 

__mul:
 # вход: buf и buf2
 # %rax - тип операции 
 # 0 - целочисленное сложение 
 # 1 - сложение вещественных чисел   
 # выход: userData 
 call __clearUserData
 cmp $0, %rax 
 jz __mulInt 
 cmp $1, %rax 
 jz __mulFloat

 call __throughError

 __mulInt:
 call __toNumber
 mov %rax, %rbx 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 call __toNumber
 imul %rbx, %rax
 call __toStr 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 ret 
 __mulFloat:
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, (buf)
 fld (buf)
 movss %xmm1, (buf)
 fmul (buf)
 fstp (buf)
 call __floatToStr

 ret 

__div:
 # вход: buf и buf2
 # только для вещественных чисел!   
 # выход: userData 
 call __clearUserData
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 # проверка деления на нуль
 movss (zero), %xmm2 
 cmpss $0, %xmm0, %xmm2
 pextrb $3, %xmm2, %rax
 cmp $0, %rax 
 jnz __divIsZero 
 movss %xmm1, (buf)
 fld (buf)
 movss %xmm0, (buf)
 fdiv (buf)
 fstp (buf)
 call __floatToStr
 ret 
 __divIsZero:
 mov $divZeroError, %rsi 
 call __throughUserError

__divI:
 # вход: buf и buf2
 # только для неотрицательных целых чисел!   
 # выход: userData 
 call __clearUserData
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set 
 call __toNumber 
 cmp $0, %rax 
 jl __divINeg
 mov %rax, %r10 # первый операнд сохранен в %r10
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __toNumber # второй операнд в %rax 
 # проверка деления на нуль  
 cmp $0, %rax 
 jz __divIsZeroI 
 jl __divINeg
 mov %rax, %rcx # запоминаем второй операнд в %rcx 
 mov %r10, %rax # первый операнд в %rax
 xor %rdx, %rdx  
 div %rcx 
 call __toStr 
 mov $lenUserData, %rsi 
 mov $userData, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 ret 
 __divIsZeroI:
 mov $divZeroError, %rsi 
 call __throughUserError
__divINeg:
 mov $divINegError, %rsi 
 call __throughUserError

__pow:
 # вход: buf и buf2
 # только для вещественных чисел!   
 # выход: userData 
 movb $0, (isExpNeg) # признак неотрицательного результата
 call __clearUserData
 call __clearBuf4
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi 
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 call __clearBuf
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi 
 call __set
 call __parseFloat
 movss %xmm1, (buf)
 movss %xmm0, (buf4)
 # основание - нуль? 
 movss (zero), %xmm2 
 movss %xmm1, %xmm3 
 cmpss $0, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax 
 jnz __powZeroBase 

 movss (zero), %xmm2 
 movss (buf), %xmm3 
 cmpss $1, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax 
 jz __powIsPos
 movss (buf4), %xmm2 
 roundps $3, %xmm2, %xmm2
 movss %xmm2, (buf3)
 fld (buf4)
 fsub (buf3)
 fstp (buf3)

 movss (zero), %xmm2 
 movss (buf3), %xmm3 
 cmpss $0, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax 
 jz __powNotInt
 cvtss2si (buf4), %rax 
 mov $2, %rbx 
 xor %rdx, %rdx 
 div %rbx 
 cmp $0, %dl # число четное?
 jnz __powNotOdd
 movb $0, (isExpNeg) # признак неотрицательного результата
 jmp __powNotOddEnd
 __powNotOdd:   
 movb $1, (isExpNeg) # признак отрицательного результата
 __powNotOddEnd:
 jmp __powInt 
 __powNotInt:
 mov $powNegError, %rsi 
 call __throughUserError
 __powInt:
 __powIsPos: 

 fldln2
 fld (buf)
 fabs 
 movss %xmm0, (buf)
 fyl2x
 fmul (buf)
 fstp (buf)
 # возводим e^buf 
 fldl2e
 fmul (buf)
 fstp (buf)
 fld (buf) # смешанное число 
 frndint
 fstp (buf2) # целое число
 fld (buf)
 fsub (buf2)
 f2xm1
 fstp (buf)
 fld1 
 fadd (buf)
 fstp (buf)
 fld (buf2)
 fld1 
 fscale 
 fmul (buf)
 fstp (buf)
 fstp (buf2)
 
 movb (isExpNeg), %al  
 cmp $0, %al    
 jz __powEnd 
 # результат отрицательный 
 fld (zero)
 fsub (buf)
 fstp (buf)
 __powEnd:
 call __floatToStr
 ret 
 __powZeroBase:
 movss (zero), %xmm2 
 movss %xmm0, %xmm3 
 cmpss $0, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax 
 jz __powZeroExpEnd  
 mov $powZeroZeroError, %rsi 
 call __throughUserError
 __powZeroExpEnd: 
 movss (zero), %xmm2 
 movss %xmm0, %xmm3 
 cmpss $1, %xmm2, %xmm3
 pextrb $3, %xmm3, %rax
 cmp $0, %rax
 jz __powNegExpEnd
 mov $powZeroNegError, %rsi 
 call __throughUserError
 __powNegExpEnd:
 call __floatToStr
 ret 

__floatToStr:
# вход: buf
# выход userData
# проверяем, является ли число отрицательным
movss (zero), %xmm0 
movss (buf), %xmm1 
cmpss $1, %xmm0, %xmm1
pextrb $3, %xmm1, %rax
cmp $0, %rax 
jz __floatToStrIsPos
# меняем знак 
fld (zero) 
fsub (one)
fmul (buf)
fstp (buf)
movb $1, (isNeg) # признак отрицательного числа 
jmp __floatToStrIsNeg
__floatToStrIsPos:
movb $0, (isNeg)
__floatToStrIsNeg:
fld (buf)
movss (buf), %xmm0 
roundps $3, %xmm0, %xmm0
movss %xmm0, (buf)
cvtss2si (buf), %r12
fsub (buf) # вычитаем из значения целое значение, получаем дробное 
fstp (buf)


mov $6, %r10 # 6 знаков после запятой 

__floatToStrLocal:
fld (buf)
cmp $0, %r10
jz __floatToStrOk
dec %r10 
movss (ten), %xmm0
movss %xmm0, (buf)
fmul (buf)
fstp (buf)
jmp __floatToStrLocal
__floatToStrOk:
cvtss2si (buf), %rax # здесь содержится дробное значение 
call __toStr

call __clearBuf
mov $buf, %rax 
movb $48, (%rax)
inc %rax 
movb $0, (%rax)

__floatToStrZeros:
mov $buf2, %rsi 
call __len 
cmp $6, %rax 
jz __floatToStrEndZeros

mov $lenBuf, %r8 
mov $buf, %r9 
mov $buf2, %r11
call __concatinate

mov $lenBuf2, %rsi 
mov $buf2, %rdx 
mov $lenUserData, %rax 
mov $userData, %rdi 
call __set

jmp __floatToStrZeros

__floatToStrEndZeros:

call __clearBuf3 # в buf3 содержится дробное значение в виде строки 
mov $lenBuf3, %rsi 
mov $buf3, %rdx 
mov $lenBuf2, %rax 
mov $buf2, %rdi 
call __set 

mov %r12, %rax 
call __toStr # в buf2 содержится целое значение в виде строки
call __clearBuf
mov $buf, %rax 
movb $'.', (%rax)
inc %rax 
movb $0, (%rax) 

mov $lenBuf2, %r8 
mov $buf2, %r9 
mov $buf, %r11 
call __concatinate



call __clearBuf 
mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenUserData, %rax 
mov $userData, %rdi 
call __set


mov $lenBuf, %r8
mov $buf, %r9 
mov $buf3, %r11 
call __concatinate 

# число отрицательное?
mov (isNeg), %dl 
cmp $1, %dl  
jnz __floatToStrEnd 
call __clearBuf
mov $buf, %rax 
movb $'-', (%rax)
inc %rax 
movb $0, (%rax)

mov $lenBuf3, %rsi 
mov $buf3, %rdx 
mov $lenUserData, %rax 
mov $userData, %rdi 
call __set 

mov $lenBuf, %r8 
mov $buf, %r9 
mov $buf3, %r11 
call __concatinate

__floatToStrEnd:
fstp (buf)
ret 

__parseFloat:
# $buf - источник (строка)
# %xmm0 - результат
mov $buf, %rax 
call __clearBuf2
call __clearBuf3
mov $buf2, %rbx # здесь будет содержаться целая часть 
mov $buf3, %rcx # здесь будет содержаться дробная часть
mov (%rax), %dl 
cmp $'-', %dl 
jnz isPos
mov $1, %r12 # признак отрицательного числа 
inc %rax 
jmp __parseFloatLocal 
isPos:
mov $0, %r12 
__parseFloatLocal: 
mov (%rax), %dl 
cmp $'.', %dl
jz __point
mov %dl, (%rbx)
inc %rax 
inc %rbx 
jmp __parseFloatLocal
__point:
movb $0, (%rbx)
mov %rax, %rbx 
mov $buf2, %rsi 
call __len 
cmp $8, %rax # целое число - не более 7 цифр 
jl __parseFloatZ
mov $buf2, %rbx 
movb $48, (%rbx)
inc %rbx 
movb $0, (%rbx) 
mov $buf3, %rbx 
movb $48, (%rbx)
inc %rbx 
movb $0, (%rbx)
jmp __parseNow 
__parseFloatZ: 
mov %rbx, %rax 

__pointLocal:             
inc %rax
mov (%rax), %dl 
cmp $0, %dl 
jz __parseNow
mov (%rax), %dl 
mov %dl, (%rcx) 
inc %rcx 
jmp __pointLocal   
__parseNow:
movb $0, (%rcx)

call __clearBuf
mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenBuf2, %rax 
mov $buf2, %rdi 
call __set 
call __toNumber
mov %rax, %r10 # целая часть числа в %r10

call __clearBuf
mov $lenBuf, %rsi 
mov $buf, %rdx 
mov $lenBuf3, %rax 
mov $buf3, %rdi 
call __set

mov $buf, %rsi 
call __len 
mov %rax, %rbx # длина дробной части числа в %rbx 
cmp $6, %rbx # дробная часть не более шести знаков 
jl __parseFloatCut
mov $buf, %rsi 
add $6, %rsi 
movb $0, (%rsi)
mov $6, %rbx 
__parseFloatCut: 
call __toNumber
mov %rax, (buf)
cvtsi2ss (buf), %xmm0  
movss %xmm0, (buf)

__floatLocal:
fld (buf)
cmp $0, %rbx 
jz __floatOk
dec %rbx 
movss (ten), %xmm0
movss %xmm0, (buf)
fdiv (buf)
fstp (buf)
jmp __floatLocal
__floatOk:
mov %r10, (buf)
cvtsi2ss (buf), %xmm0    
movss %xmm0, (buf) # целая часть числа 
fadd (buf)
fstp (buf)
movss (buf), %xmm0  
cmp $1, %r12 
jnz __pos
mov (zero), %rax   
mov %rax, (buf)
fld (buf)
movss %xmm0, (buf)
fsub (buf)
fstp (buf)
movss (buf), %xmm0 
__pos:
ret 

__initLabels:
 # формирует в %rax адрес начала выделяемой памяти для firstMem
 mov $12, %rax
 xor %rdi, %rdi
 syscall
 mov %rax, (memoryBegin)
 call __newLabelMem
 mov (memoryBegin), %rax 
 add (labelSize), %rax 
 mov %rax, %r12 # сохраняем rax 
 mov (memoryBegin), %rdi
 mov (memoryBegin), %r9
 
 
 mov %rdi, %r10 # сохраняем %rdi  
 mov %r9, %rsi # сохраняем %rsi 
 mov %r12, %rax 
 
 call __newLabelMem
 add (labelSize), %r12 

 mov %r10, %rdi # восстанавливаем  
 mov %rsi, %r9 
 
 mov $labelName1, %rbx
 __initLabelsName1: 
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsNameEx1
 mov %dl, (%rdi) 
 inc %rbx 
 inc %rdi 
 jmp __initLabelsName1 
 __initLabelsNameEx1:
 movb $0, (%rdi)

 mov (label1), %rax 
 call __toStr
 add (valSize), %r9
 mov %r9, %rdi 
 mov $buf2, %rbx 
__initLabelsAddr1:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx1
 mov %dl, (%rdi)
 inc %rbx 
 inc %rdi 
 jmp __initLabelsAddr1
 __initLabelsAddrEx1:
 movb $0, (%rdi)

 add (valSize), %r9 
 mov %r9, %rdi 
 

 mov %rdi, %r10 # сохраняем %rdi  
 mov %r9, %rsi # сохраняем %rsi 
 mov %r12, %rax
 
 call __newLabelMem
 add (labelSize), %r12 

 mov %r10, %rdi # восстанавливаем  
 mov %rsi, %r9
  
 mov $labelName2, %rbx 
 
 __initLabelsName2:
 mov (%rbx), %dl 
 cmp $0, %dl  
 jz __initLabelsNameEx2
 mov %dl, (%rdi)
 inc %rbx 
 inc %rdi 
 jmp __initLabelsName2
 __initLabelsNameEx2:
 movb $0, (%rdi)
 mov (label2), %rax 
 call __toStr
 add (valSize), %r9 
 mov %r9, %rdi 
 mov $buf2, %rbx 
 __initLabelsAddr2:
 mov (%rbx), %dl 
 cmp $0, %dl 
 jz __initLabelsAddrEx2
 mov %dl, (%rdi)
 inc %rbx 
 inc %rdi 
 jmp __initLabelsAddr2
 __initLabelsAddrEx2:
 movb $0, (%rdi)

 mov %r12, %rax 
 mov %r12, (labelsMax)
 ret 

 __goto:
 # адрес имени метки, по которой нужно прыгнуть, в %rdi 
 mov %rdi, %rsi 
 call __len 

 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 call __set 

 mov (memoryBegin), %rax
  __gotoSearch: 
 mov %rax, (labelsPointer)
 cmp %rax, (labelsMax)
 jl __gotoEnd
 mov %rax, %r12 
 call __read  
 
 call __compare 
 cmp $1, %rax 
 jz __goNow  
 mov (labelsPointer), %rax 
 add (labelSize), %rax 
  
 jmp __gotoSearch
 __goNow:
 mov (labelsPointer), %rax 
 add (valSize), %rax 
 mov %rax, %r8 #  cохраняем %rax 
 mov %rax, %rsi 
 call __len 
 
  
 mov $lenBuf, %rsi 
 mov $buf, %rdx  
 mov %r8, %rdi 
 call __set 
 
 call __toNumber
 jmp *%rax 

 __gotoEnd: 
 mov $noSuchMarkError, %rsi 
 call __len 
 mov %rax, %r8 
 mov $noSuchMarkError, %r9 
 mov $buf2, %r11 
 call __concatinate

 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenUserData, %rax 
 mov $userData, %rdi
 call __set 
 
 mov $lenUserData, %r9 
 mov $buf, %r9 
 mov $enter, %r11 
 call __concatinate

 mov $userData, %rsi 
 call __throughUserError
 ret 

 __less:
 # вход: buf и buf2 
 # %rax - тип операции 
 # 0 - целочисленный
 # 1 - вещественный 
 # выход: userData 
 cmp $0, %rax 
 jnz __lessFloat
 call __toNumber
 mov %rax, %r12 # сохраняем %rax 
 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set 
 call __toNumber 
 
 cmp %r12, %rax 
 jg __isLess 
 call __clearUserData
 movb $'0', (userData)
 ret 
 __isLess:
 call __clearUserData
 movb $'1', (userData)
 ret 
 __lessFloat:
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __parseFloat

 cmpss $2, %xmm1, %xmm0
 pextrb $3, %xmm0, %rax
 cmp $0, %rax 
 jz __isLess 
 call __clearUserData
 movb $'0', (userData)
 ret 

  __lessOrEqual:
 # вход: buf и buf2 
 # %rax - тип операции 
 # 0 - целочисленный
 # 1 - вещественный 
 # выход: userData 
 cmp $0, %rax 
 jnz __lessOrEqualFloat
 call __toNumber
 mov %rax, %r12 # сохраняем %rax 
 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set 
 call __toNumber 
 
 cmp %r12, %rax 
 jge __isLessOrEqual 
 call __clearUserData
 movb $'0', (userData)
 ret 
 __isLessOrEqual:
 call __clearUserData
 movb $'1', (userData)
 ret 
 __lessOrEqualFloat:
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __parseFloat

 cmpss $1, %xmm1, %xmm0
 pextrb $3, %xmm0, %rax
 cmp $0, %rax 
 jz __isLessOrEqual 
 call __clearUserData
 movb $'0', (userData)
 ret 

 __more:
 call __lessOrEqual
 xor %rax, %rax 
 mov (userData), %al
 cmp $'0', %al 
 jz __isMore
 movb $'0', (userData)
 ret 
 __isMore:
 movb $'1', (userData)
 ret 

 __moreOrEqual:
 call __less
 xor %rax, %rax 
 mov (userData), %al
 cmp $'0', %al 
 jz __isMoreOrEqual
 movb $'0', (userData)
 ret 
 __isMoreOrEqual:
 movb $'1', (userData)
 ret

 __eq:
 # вход: buf и buf2 
 # %rax - тип операции 
 # 0 - целочисленный
 # 1 - вещественный 
 # выход: userData 
 cmp $0, %rax 
 jnz __equalFloat
 call __toNumber
 mov %rax, %r12 # сохраняем %rax 
 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set 
 call __toNumber 
 
 cmp %r12, %rax 
 jz __isEqual  
 call __clearUserData
 movb $'0', (userData)
 ret 
 __isEqual:
 call __clearUserData
 movb $'1', (userData)
 ret 
 __equalFloat:
 mov $lenBuf4, %rsi 
 mov $buf4, %rdx 
 mov $lenBuf2, %rax 
 mov $buf2, %rdi
 call __set
 call __parseFloat
 movss %xmm0, %xmm1 
 mov $lenBuf, %rsi 
 mov $buf, %rdx 
 mov $lenBuf4, %rax 
 mov $buf4, %rdi
 call __set
 call __parseFloat

 cmpss $4, %xmm1, %xmm0
 pextrb $3, %xmm0, %rax
 cmp $0, %rax 
 jz __isEqual  
 call __clearUserData
 movb $'0', (userData)
 ret 

 __parseBool:
 # buf - источник (строка)
 # %rax - результат

 xor %rax, %rax 
 mov (buf), %al  
 cmp $'1', %al 
 jnz __parseFalse
 mov $1, %rax 
 ret  
 __parseFalse:
 mov $0, %rax 
 ret 

 __boolToStr:
 # вход: buf
 # выход: userData
 call __clearUserData
 mov (buf), %al 
 cmp $1, %al 
 jnz __boolToStrEndTrue
 movb $'1', (userData)
 ret 
 __boolToStrEndTrue:
 movb $'0', (userData)
 ret 

 __and:
 # вход: buf и buf2 в виде строк 
 # выход: userData в виде строки 
 call __clearUserData
 call __parseBool 
 mov %rax, (userData)
 mov (buf2), %rax 
 mov %rax, (buf)
 call __parseBool
 mov %rax, (buf2)
 mov (userData), %rax 
 mov %rax, (buf)


 mov (buf), %rax 
 and (buf2), %rax  

 cmp $1, %rax 
 jz __andTrue 
 movb $'0', (userData)
 ret 
 __andTrue:
 movb $'1', (userData)
 
 ret 

 __or:
 # вход: buf и buf2 в виде строк 
 # выход: userData в виде строки 
 call __clearUserData
 call __parseBool 
 mov %rax, (userData)
 mov (buf2), %rax 
 mov %rax, (buf)
 call __parseBool
 mov %rax, (buf2)
 mov (userData), %rax 
 mov %rax, (buf)


 mov (buf), %rax 
 or (buf2), %rax  

 cmp $1, %rax 
 jz __orTrue 
 movb $'0', (userData)
 ret 
 __orTrue:
 movb $'1', (userData)
 
 ret 

 __xor:
 # вход: buf и buf2 в виде строк 
 # выход: userData в виде строки 
 call __clearUserData
 call __parseBool 
 mov %rax, (userData)
 mov (buf2), %rax 
 mov %rax, (buf)
 call __parseBool
 mov %rax, (buf2)
 mov (userData), %rax 
 mov %rax, (buf)


 mov (buf), %rax 
 xor (buf2), %rax  

 cmp $1, %rax 
 jz __xorTrue 
 movb $'0', (userData)
 ret 
 __xorTrue:
 movb $'1', (userData)
 
 ret 

__not:
 # вход: buf в виде строки 
 # выход: userData в виде строки 
 call __clearUserData

 mov (buf), %al 
 cmp $'1', %al 
 jz __notTrue
 movb $'1', (userData)
 ret 
 __notTrue:
 movb $'0', (userData)
 ret 

.globl _start
_start:
 call __initLabels
 call __firstMem
 call __firstStrMem

 # ^systemVar 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenSystemVarName, %rax 
 mov $systemVarName, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar 

 # sVar 
 /*mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1, %rax 
 mov $varName1, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar*/

 # iVar 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName0, %rax 
 mov $varName0, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi
 call __set 
 call __defineVar

 

  # bVar 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName8, %rax 
 mov $varName8, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenBoolType, %rax 
 mov $boolType, %rdi
 call __set 
 call __defineVar

# iVar2  
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
 mov $varName4, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenIntType, %rax 
 mov $intType, %rdi
 call __set 
 call __defineVar
 
 #set iVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName0, %rax 
 mov $varName0, %rdi
 call __set 
 mov $data3, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar

 #set bVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName8, %rax 
 mov $varName8, %rdi
 call __set 
 mov $data10, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar

 #set iVar2
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
 mov $varName4, %rdi
 call __set 
 mov $data0, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
 
 # fVar 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName2, %rax 
 mov $varName2, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenFloatType, %rax 
 mov $floatType, %rdi
 call __set 
 call __defineVar

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName2, %rax 
 mov $varName2, %rdi
 call __set 
 mov $data4, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar

  # fVar2
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName5, %rax 
 mov $varName5, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenFloatType, %rax 
 mov $floatType, %rdi
 call __set 
 call __defineVar

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName5, %rax 
 mov $varName5, %rdi
 call __set 
 mov $data5, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
 
 # sVar 
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1, %rax 
 mov $varName1, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar


  # sVar2
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName6, %rax 
 mov $varName6, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar

  # sVar3
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName7, %rax 
 mov $varName7, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar

  # sVar2
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName9, %rax 
 mov $varName9, %rdi
 call __set 
 mov $lenVarType, %rsi 
 mov $varType, %rdx 
 mov $lenStringType, %rax 
 mov $stringType, %rdi
 call __set 
 call __defineVar
  
  #set sVar
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName1, %rax 
 mov $varName1, %rdi
 call __set 
 mov $data1, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar

 #set sVar2
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName6, %rax 
 mov $varName6, %rdi
 call __set 
 mov $data2, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar
  
  
  #set sVar3
 /*mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName7, %rax 
 mov $varName7, %rdi
 call __set 
 mov $data8, %rax 
 mov %rax, (userData)
 xor %rax, %rax 
 call __setVar

 # set sVar4
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName9, %rax 
 mov $varName9, %rdi
 call __set 
 mov $varName7, %rax 
 mov %rax, (userData)
 mov $1, %rax 
 call __setVar */
 

 # get fVar  
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName2, %rax 
 mov $varName2, %rdi
 call __set
 call __getVar

 mov (userData), %rsi 
 call __len 
 mov $lenBuf3, %rsi 
 mov $buf3, %rdx 
 mov (userData), %rdi 
 call __set 

 # get fVar2
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName5, %rax 
 mov $varName5, %rdi
 call __set
 call __getVar

 mov (userData), %rsi 
 call __len 
 mov $lenBuf2, %rsi 
 mov $buf2, %rdx 
 mov (userData), %rdi 
 call __set 

 mov $lenBuf, %rsi 
 mov $buf, %rdx
 mov $lenBuf3, %rax 
 mov $buf3, %rdi 
 call __set

 call __clearBuf
 movb $'1', (buf)
 call __not  



  # get fVar  
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName2, %rax 
 mov $varName2, %rdi
 call __set
 call __getVar

 push (userData)

 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName0, %rax 
 mov $varName0, %rdi
 call __set
 call __getVar

 push (userData)

 pop (userData)
 
  #set iVar2
 mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName4, %rax 
 mov $varName4, %rdi
 call __set 
 xor %rax, %rax 
 call __setVar

 # подготавливаем вызов __userConcatinate

 mov $lenVarName, %rsi 
 mov $varName, %rdx
 mov $lenVarName1, %rax 
 mov $varName1, %rdi 
 call __set

 #mov (userMem), %r8 
 #mov (userMem2), %r9
 
 mov $varName1, %r8
 mov $data12, %r9 
 mov $1, %rax 
 mov $0, %rbx 

 /*mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName6, %rax 
 mov $varName6, %rdi
 call __set
 call __undefineVar
 call __printHeap*/

 call __userConcatinate 

  

 #get sVar 
 /*mov $lenVarName, %rsi 
 mov $varName, %rdx 
 mov $lenVarName6, %rax 
 mov $varName6, %rdi
 call __set
 call __getVar
 mov (userData), %rsi 
 call __print 
 mov $enter, %rsi 
 call __print*/ 

 call __printHeap

__stop:
 mov $60,  %rax      # номер системного вызова exit
 xor %rdi, %rdi      # код возврата (0 - выход без ошибок)
 syscall   

__throughMessage:
 mov $data1, %rsi 
 call __print 
 mov $60,  %rax      # номер системного вызова exit
 xor %rdi, %rdi      # код возврата (0 - выход без ошибок)
 syscall 
