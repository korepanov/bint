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

stack get_funcs(){
	stack res_stack;
	string return_type;
	string func_name;
	
	#get_funcs_s:
	func_stack = next_func();
	func_stack.pop(return_type);
	func_stack.pop(func_name);
	
	[goto(#get_funcs_e), ("end" == func_name), print("")];
	res_stack.push(func_name);
	res_stack.push(return_type); 
	goto(#get_funcs_s);
	#get_funcs_e:
	RESET_SOURCE();
	return res_stack;
		
};

stack reverse(stack s){
	string buf;
	stack res; 
	
	s.pop(buf);
	#reverse_s:
	[goto(#reverse_e), ("end" == buf), print("")];
	res.push(buf);
	s.pop(buf);
	goto(#reverse_s);
	#reverse_e:
	return res;
};

stack indexes(string s, string sub_s){
	stack res;
	
	int i;
	int pointer; 
	int s_len;
	int s_len_old;
	int sub_len;
 	
	s_len = len(s);
	sub_len = len(sub_s);
	s_len_old = s_len;
	i = index(s, sub_s);
	pointer = i;
	#indexes_s:
	[goto(#indexes_e), (-1 == i), print("")];
	i = (i + (s_len_old - s_len));
	res.push(i);
	pointer = (pointer + sub_len);
	s = s[pointer:s_len];
	s_len = len(s);
	i = index(s, sub_s);
	pointer = i;
	goto(#indexes_s);
	#indexes_e:
	res = reverse(res);
	return res;
};

void replace(){
	string command;
	int command_len;
	int number;
	int func_pos;
	stack func_pos_stack;
	bool change_flag;
	int func_len;
	string symbol;
	string return_type;
	int func_entry;
	string str_func_entry;
	string buf;

	func_entry = 0;
	change_flag = False;
	func_stack = get_funcs();
	
	#replace_s:
	func_stack.pop(return_type);
	func_stack.pop(func_name);
	
	[goto(#replace_e), ("end" == func_name), print("")];
	
	#next:
	next_command(command);
	[goto(#next_end), ("end" == command), print("")];
	number = index(command, func_name);
	[print(""), (-1 == number), goto(#not_send)];
	send_command(command);
	goto(#next);
	#not_send:
	
	func_len = len(func_name);
	number = (number + func_len);
	symbol = command[number];
	[print(""), ("(" == symbol), goto(#next)];

	
	arg_type = "int";
	number = index(command, "int");
	[goto(#next), (0 == number), print("")];
	arg_type = "bool";
	number = index(command, "bool");
	[goto(#next), (0 == number), print("")];
	arg_type = "float";
	number = index(command, "float");
	[goto(#next), (0 == number), print("")];
	arg_type = "stack";
	number = index(command, "stack");
	[goto(#next), (0 == number), print("")];
	arg_type = "string";
	number = index(command, "string");
	[goto(#next), (0 == number), print("")];
	
	
	print(command);
	print("\n");
	func_pos_stack = indexes(command, func_name);	
	func_pos_stack.pop(buf);
	#pop_func_pos_start:
	[goto(#pop_func_pos_end), ("end" == buf), print("")];
	print(buf);
	print("\n");
	func_pos_stack.pop(buf);
	goto(#pop_func_pos_start);
	#pop_func_pos_end:

	str_func_entry = str(func_entry);
	command = (((("$" + func_name) +  "_res") + str_func_entry) + "=");
	func_entry = (func_entry + 1);
	

	print(command);
	print("\n");
	
	
	send_command(command);
	goto(#next);
	#next_end:
	UNSET_SOURCE();
	UNSET_DEST();
	[goto(#change), (change_flag), print("")];
	change_flag = True;
	SET_SOURCE("benv/long_function_program.b");
	SET_DEST("benv/long_function_program2.b");
	goto(#replace_s);
	#change:
	change_flag = False; 
	SET_SOURCE("benv/long_function_program.b");
	SET_DEST("benv/long_function_program2.b");
	goto(#replace_s);

	#replace_e:
	print("");
};

int res; 

res = init();

[print(""), (0 == res), print("INIT ERROR\n")];
replace();
res = finish();
[print(""), (0 == res), print("FINISH ERROR\n")];
