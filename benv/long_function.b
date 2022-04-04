string command;
string return_type;
string func_name;
stack func_stack;
bool bool_res;

int init(){
	SET_SOURCE("benv/prep_func_program.b");
	SET_DEST("benv/long_function_program.b");
	
	return 0;
};

int finish(){
	DEL_DEST("benv/prep_func_program.b");
	UNSET_SOURCE();
	UNSET_DEST();

	return 0;
};

stack next_func(){
	int number;
	int left_border;
	int right_border;
	stack func_stack;
	string func_name;
	int type_len;
	string command;
	string arg_type;
	
	#next_func_st:
	next_command(command);
	send_command(command);
	[goto(#end_file), ("end" == command), print("")];
	number = index(command, "{");
	[goto(#end_clause), (-1 == number), print("")];
	arg_type = "int";
	number = index(command, "int");
	[goto(#next_func_e), (0 == number), print("")];
	arg_type = "bool";
	number = index(command, "bool");
	[goto(#next_func_e), (0 == number), print("")];
	arg_type = "float";
	number = index(command, "float");
	[goto(#next_func_e), (0 == number), print("")];
	arg_type = "stack";
	number = index(command, "stack");
	[goto(#next_func_e), (0 == number), print("")];
	arg_type = "string";
	number = index(command, "string");
	[goto(#next_func_e), (0 == number), print("")];
	#end_clause:
	goto(#next_func_st);
	#next_func_e:
	type_len = len(arg_type);
	left_border = type_len; 
	right_border = index(command, "(");
	
	func_name = command[left_border:right_border];
	
	func_stack.push(func_name);
	func_stack.push(arg_type);
	
	#end_file:
	return func_stack;
};

int res; 
res = init();

[print(""), (0 == res), print("INIT ERROR\n")];
func_stack = next_func();
func_stack.pop(return_type);
func_stack.pop(func_name);
print(return_type);
print(" ");
print(func_name);
print("\n");
func_stack = next_func();
func_stack.pop(return_type);
func_stack.pop(func_name);
print(return_type);
print(" ");
print(func_name);
print("\n");
res = finish();
[print(""), (0 == res), print("FINISH ERROR\n")];
