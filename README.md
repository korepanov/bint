# Интерпретатор языков B и Basm 
## Синтаксис языка Basm
Язык **Basm** - **Assembler**-подобный язык

Каждая команда заканчивается символом "**;**"\
Типы данных - **float, int, string, bool, stack**\
Возможно явное приведение типов данных функциями **float(), int(), str(), bool()**.
### Пример явного преобрзования типов 
```
string buf;
int a;
a = int(buf);
```
### Синтаксис условной дизъюнкции
```
[goto(#mark1)/print("text1"), (<условие>), goto(#mark2)/print("text2")];
```
В случае, когда условие истинно, выполняется команда слева от условия (**goto(#mark1)**, либо **print("text1")**).
В противном случае выполняется команда справа.\
В условии обязательны скобки для каждой элементарной операции,
а также для каждого операнда перед логической операцией (см. ниже)

Логические операторы: **AND, OR, XOR, NOT**
### Пример логического выражения
```
bool b
b = ((True)AND(False));
```
Метка объявляется следующим образом: 
### Пример метки 
```
#mark1:
```
Оператор goto позволяет перейти к метке, объявленной в любом месте программы:
### Пример использования оператора goto
```
goto(#mark);
```
### Условная дизъюнкция 
```
[goto(#left), (4 > 5), goto(#right];
```
В данном случае интерпретатор перейдет к метке **#right**\
Как было сказано выше, скобки обязательны для каждой элементарной операции.
Приведем пример корректного арифметического выражения 
### Корректное арифметическое выражение 
```
float res; 
res = (((5^2) + 4) - 3);
```
Арифметические операции: **+**, **-**, __*__, **/**, **^**\
Оператор **print** печатает выражение на экран 
### Пример использования print 
```
print("Hello world!\n");
```
*Примечание*: русский язык не поддерживается\
Оператор **len** позволяет вычислить длину строки 
### Пример использования len 
```
int a; 
string buf; 
buf = "Hello world!";
a = len(buf);
```
*Примечание*: использование len возможно только после переменной и знака присваивания\
Оператор **index** позволяет определить индекс первого вхождения подстроки в строку 
### Синтаксис index 
```
index(<строка>, <подстрока>);
```
### Пример использования index 
```
int a;
a = index("banana", "nan");
```
В случае отсуствия вхождения **index** возвращает -1.\
Как и для оператора **len**, использование **index** возможно только после переменной и знака присваивания\
Оператор **is_letter** возвращает True в случае, когда аргумент является буквкой и False в противном случае
### Синтаксис is_letter
```
is_letter(<символ>)
```
### Пример использования is_letter
```
bool a;
a = is_letter("b");
string buf;
buf = str(a);
print(buf);
print("\n");
```
Аналогично работает оператор **is_digit**, возвращающий значение True в случае, 
когда символ является цифрой и возвращающий False в противном случае. 
Оператор **push** позволяет положить переменную любого типа в стек
### Пример использования push 
```
push(100);
```
Оператор **pop** позволяет достать содержимое с вершины стека.\
При этом переменная соответствующего типа передается оператору pop в качестве аргумента 
### Пример использования pop для того, чтобы достать с вершины стека число 100
```
int a; 
pop(a);
``` 
Теперь в переменной a находится число 100 
Специальный тип stack позволяет определять стек пользовательского типа.\
Чтобы применять к переменной типа stack операции push и pop необходимо после
переменной поставить символ "**.**"
### Пример использования пользовательского стека
```
stack my_stack; 
my_stack.push(5);
my_stack.push("Hello world!"); 
string buf; 
my_stack.pop(buf);
print(buf);
print("\n");
int a; 
my_stack.pop(a); 
buf = str(a); 
print(buf);
print("\n");
```
Операция **SET_SOURCE("<файл>")** открывает системный файл на чтение 
### Пример использования SET_SOURCE
```
SET_SOURCE("program.b");
```
Операция **SET_DEST("<файл>")** открывает системный файл на запись\
Операция **UNSET_SOURCE()** закрывает системный файл, открытый на чтение\
Операция **UNSET_DEST()** закрывает системный файл, открытый на запись

*Примечание*: единовременно может быть открыт только один системный файл на запись и
один системный файл на чтение

Операция **next_command(<переменная типа string>)** позволяет считать очередную 
команду из системного файла до символа "**;**"\
Операция **send_command(<переменная типа string>)** позволяет послать очередную команду
на запись 
### Пример считывания и пересылки команды 
```
SET_SOURCE("program.b");
SET_DEST("program.basm");
string command; 
next_command(command);
send_command(command);
UNSET_SOURCE();
USET_DEST();
```

Оператор **UNDEFINE(<переменная>)** сообщает интерпретатору о том, что
необходимо "забыть" о существовании переменной\
Логические значения: **True** и **False**\
Операции сравнения: **<, <=, ==, >=, >**\
Примечание: операция "!=" отсутствует. Вместо этого стоит использовать
оператор **NOT** 
### Пример условия неравенства
```
int a; 
int b; 
a = 5; 
b = 10;
[goto(#left), (NOT(a == b)), goto(#right)];
```
Для ввода данных используется оператор **input**
### Пример использования input 
```
string buf; 
input(buf);
```
Поскольку модули трансляции новых конструкций языка B на язык Basm пишутся либо на языке Basm,\
либо на старых конструкциях языка B, то внутри этих модулей необходимо знать путь к файлу-источнику и\
к файлу назначения. Для определения этих путей используются операторы **get_root_source** и **get_root_dest**
### Пример использования get_root_source и get_root_dest 
```
string a;
get_root_source(a);
print(a);
print("\n");
get_root_dest(a);
print(a);
print("\n");
```
Операция **SEND_DEST("<имя файла>")** пересылает файл с именем "program.basm"\
(траслированный с языка B файл) из папки benv по пути назначения, присваивая ему имя
<имя файла>
### Пример использования SEND_DEST
```
SEND_DEST("res.basm");
```
Операция **DEL_DEST("<файл>")** удаляет файл по пути, указанному в качестве аргумента, в случае,
когда этот файл существует
### Пример использования DEL_DEST 
```
DEL_DEST("benv/res.b");
```
Оператор **REROUTE()** меняет местами файл-источник и файл назначения 
### Пример использования REROUTE
```
REROUTE();
```
Оператор **RESET_SOURCE()** сбрасывает счетчик команд файла-источника и позволяет прочесть файл
сначала
### Пример использования RESET_SOURCE()
```
RESET_SOURCE()
```

Пробелы, символы табуляции и перехода на следующую строку
игнорируются\
Комментарии запрещены\
Приведем пример программы на языке Basm
### Пример программы решения квадратных уравнений на языке Basm
```
print("Solving equation of the form ax^2 + bx + c = 0\n");
#begin:
print("Input a\n");
float a;
string buf;
input(buf);
a = float(buf);
print("Input b\n");
float b;
input(buf);
b = float(buf);
print("Input c\n");
float c;
input(buf);
c = float(buf);
float d;
d = ((b^2) - ((4*a)*c));
float x1;
float x2;
[print(""), (d >= 0), goto(#no_solution)];
x1 = ((((-1)*b) - (d^0.5)) / (2*a));
x2 = ((((-1)*b) + (d^0.5)) / (2*a));
print("x1 = ");
buf = str(x1);
print(buf);
print("\n");
print("x2 = ");
buf = str(x2);
print(buf);
print("\n");
goto(#end_iter);
#no_solution:
print("No solution\n");
#end_iter:
print("Exit? y/n\n");
input(buf);
[print(""), ("y" == buf), goto(#begin)];
```
## Синтаксис языка B 
**B** является **C**-подобным языком над языком **Basm**. Язык **B** транслируется в язык **Basm**, а 
затем полученная программа интерпретируется. Язык **B** находится в разработке. В настоящее время
язык **B** является языком **Basm** с возможностью использования **C**-подобных функций. Метки обладают 
глобальной областью видимости, что важно учитывать при разработке 
### Пример программы решения квадратных уравнений на языке B
```
void println(string s){
	print(s);
	print("\n");
};

float descr(float a, float b, float c){
	return ((b^2) - ((4*a)*c));
};

stack solve(float a, float b, float c){
	float x1;
	float x2;
	stack res;
	float d; 
	d = descr(a, b, c);
	[goto(#solve_e), (d < 0), print("")];
	x1 = ((((-1)*b) - (d^0.5)) / (2 * a));
	x2 = ((((-1)*b) + (d^0.5)) / (2 * a));
	res.push(x2);
	res.push(x1);
	#solve_e:
	return res;
};

stack input_data(){
	float a;
	float b;
	float c;
	string buf;
	
	println("Input a");
	input(buf);
	a = float(buf);
	
	println("Input b");
	input(buf);
	b = float(buf);
	
	println("Input c");
	input(buf);
	c = float(buf);
	
	stack res;
	res.push(c);
	res.push(b);
	res.push(a);
	return res;
		
};

void output(float x1, float x2){
	string buf;
	println("Solution:");
	print("x1=");
	buf = str(x1);
	println(buf);
	print("x2=");
	buf = str(x2);
	println(buf);
};
void main(){
	println("Solving equation of the form ax^2 + bx + c = 0");
	stack res;
	res = input_data();
	float a;
	float b;
	float c; 
	res.pop(a);
	res.pop(b);
	res.pop(c);
	res = solve(a, b, c);
	string buf;
	res.pop(buf);
	[goto(#no_solution), ("end" == buf), print("")];
	res.push(buf);
	float x1;
	float x2;
	res.pop(x1);
	res.pop(x2);
	output(x1, x2);
	goto(#main_e);
	#no_solution:
	println("No solution");
	#main_e:
	print("");
};

main();
```
## Начать работу 
### Трансляция программы с языка B на язык Basm 
```
./bint -i input.b -o output.basm
```
### Интерпретировать программу на языке Basm 
```
./bint -e output.basm 
```
## Примечания 
Файл **b.lang** описывает синтаксис языков **B** и **Basm** для текстового редактора gedit (протестировано
на Ubuntu 18.04).
Чтобы gedit распознал языки **B** и **Basm**, необходимо добавить данный файл по следующему пути:
```
usr/share/gtksourceview-3.0/language-specs
```