string command;
string return_type;
string func_name;
stack func_stack;
bool bool_res;
string root_source;

int init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/long_function_program.b");
	
	return 0;
};

int finish(){
	
	return 0;
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

stack look_behind(string reg, string s){
	stack st;
	stack this;
	stack res;  
	string buf;
	int pos;
	string symbol;

	st = reg_find(reg, s);

	#look_behind_s:
	st.pop(this);
	this.pop(buf);
	[goto(#look_behind_e), ("end" == buf), print("")];
	pos = int(buf);
	[print(""), (0 == pos), goto(#is_not_zero)];
	res.push("$");
	goto(#look_behind_s);
	#is_not_zero:
	pos = (pos - 1);
	symbol = s[pos];
	res.push(symbol);
	goto(#look_behind_s);
	#look_behind_e:
	return reverse(res);
};

int func_call_index(string command, string func_name){
	string reg;
	stack st;
	stack ist;
	int res;
	string buf;
	string symbol;
	bool letter;
	bool digit;

	reg = (("(?:" + func_name) + ")");

	st = look_behind(reg, command);
	ist = indexes(command, func_name);
	
	#func_call_index_s:
	st.pop(symbol);
	
	if ("end" == symbol){
		return -1;	
	};

	ist.pop(buf);

	if ("end" == buf){
		print("func_call_index ERROR\n");	
	}; 

	res = int(buf);
	letter = is_letter(symbol);
	digit = is_digit(symbol);	

	if (NOT(((letter)OR(digit))OR("_" == symbol))){
		return res;		
	};

	goto(#func_call_index_s);
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

stack func_ends(string command, stack func_begins, int func_len){
	string buf;
	string symbol;
	int i;
	int br_begin;
	int br_end;
	int command_len;
	int opened_braces;
	int closed_braces;
	stack res;
	string temp;
	
	
	func_begins.pop(buf);
	#func_ends_s:
	closed_braces = 0;
	[goto(#func_ends_e), ("end" == buf), print("")];
	i = int(buf);
	command_len = len(command);
	br_begin = (i + func_len);
	opened_braces = 1;
	br_end = (br_begin + 1);
	#counter_s:
	[goto(#counter_e), (opened_braces == closed_braces), print("")];
	symbol = command[br_end];
	[print(""), ("(" == symbol), goto(#inc_o_end)];
	opened_braces = (opened_braces + 1);
	#inc_o_end:
	[print(""), (")" == symbol), goto(#inc_c_end)];
	closed_braces = (closed_braces + 1);
	#inc_c_end:
	br_end = (br_end + 1);
	
	goto(#counter_s);	
	#counter_e:
	res.push(br_end);
	func_begins.pop(buf);
	goto(#func_ends_s);
	#func_ends_e:
	return reverse(res);
};

void del_file(bool change_flag){
	string command;
	[print(""), (change_flag), goto(#copy_e)];
	SET_SOURCE("benv/long_function_program2.b");
	SET_DEST("benv/long_function_program.b");
	next_command(command);
	#copy_s:
	[goto(#copy_e), ("end" == command), print("")];
	send_command(command);
	next_command(command);
	goto(#copy_s);
	#copy_e:
	DEL_DEST("benv/long_function_program2.b");
};

void replace(){
	string command;
	string command_to_send;
	string replaced_command;
	string left_part;
	string right_part;
	int command_len;
	int number;
	int func_pos;
	int itemp;
	int offset;
	stack func_pos_stack;
	stack func_ends_stack;
	bool change_flag;
	int func_len;
	string symbol;
	string return_type;
	int func_entry;
	string str_func_entry;
	string sleft_border;
	string sright_border;
	int left_border;
	int right_border;
	int left_border_reserv;
	int right_border_reserv;
	string func_call;
	string stemp;
	int itemp;
	string temp;
	int stemp_len;
	string debug_var;

	func_entry = 0;
	offset = 0;
	change_flag = False;
	was_replace = False;
	func_stack = get_funcs();
	func_stack.push("$temp");
	func_stack.push("$temp");

	#replace_s:
	func_stack.pop(return_type);
	func_stack.pop(func_name);
		
	[goto(#replace_e), ("end" == func_name), print("")];
	
	#next:
	next_command(command);
	[goto(#next_end), ("end" == command), print("")];
	number = index(command, func_name);
	itemp = func_call_index(command, func_name);
	debug_var = str(itemp);
	print(command);
	print("\n");
	print(func_name);
	print("\n");
	print(debug_var);
	print("\n");
	[print(""), (-1 == number), goto(#not_send)];
	send_command(command);
	goto(#next);
	#not_send:
	
	func_len = len(func_name);
	number = (number + func_len);
	symbol = command[number];
	
	[print(""), ("(" == symbol), goto(#to_next_start)];

	arg_type = "int";
	number = index(command, "int");
	[goto(#to_next_start), (0 == number), print("")];
	arg_type = "bool";
	number = index(command, "bool");
	[goto(#to_next_start), (0 == number), print("")];
	arg_type = "float";
	number = index(command, "float");
	[goto(#to_next_start), (0 == number), print("")];
	arg_type = "stack";
	number = index(command, "stack");
	[goto(#to_next_start), (0 == number), print("")];
	arg_type = "string";
	number = index(command, "string");
	[goto(#to_next_start), (0 == number), goto(#to_next_end)];
	
	#to_next_start:
	send_command(command);
	goto(#next);
	#to_next_end:
	
	func_pos_stack = indexes(command, func_name);	
	func_ends_stack = func_ends(command, func_pos_stack, func_len);

	func_pos_stack.pop(sleft_border);
	func_ends_stack.pop(sright_border);
	replaced_command = command;
	itemp = len(command);
	stemp_len = 0;

	#pop_func_pos_start:
	[goto(#pop_func_pos_end), ("end" == sleft_border), print("")];
	left_border = int(sleft_border);
	right_border = int(sright_border);
	left_border_reserv = left_border;
	right_border_reserv = right_border;

	left_border = (left_border + offset);
	right_border = (right_border + offset);

	str_func_entry = str(func_entry);
	command_to_send = ((((return_type + "$") + func_name) + "_res") + str_func_entry);
	send_command(command_to_send);

	command_to_send = (((("$" + func_name) +  "_res") + str_func_entry) + "=");
	func_entry = (func_entry + 1);
	func_call = command[left_border_reserv:right_border_reserv];
	left_part = replaced_command[0:left_border];
	right_part = replaced_command[right_border:itemp];
	
	replaced_command = (((((left_part + "$") + func_name) + "_res") + str_func_entry) + right_part); 
	stemp = ((("$" + func_name) + "_res") + str_func_entry);
	stemp_len = len(stemp);
	offset = (offset + (stemp_len - (right_border - left_border)));
	temp = str(offset);
	itemp = len(replaced_command);
	command_to_send = (command_to_send + func_call);
	
	send_command(command_to_send);
	
	func_pos_stack.pop(sleft_border);
	func_ends_stack.pop(sright_border);
	goto(#pop_func_pos_start);
	#pop_func_pos_end:
	send_command(replaced_command);
	offset = 0;
	goto(#next);
	#next_end:
	UNSET_SOURCE();
	UNSET_DEST();
	func_stack.pop(func_name);
	func_entry = 0;
	offset = 0;
	[goto(#replace_e), ("end" == func_name), print("")];
	func_stack.push(func_name);
	[goto(#change), (change_flag), print("")];
	change_flag = True;
	SET_SOURCE("benv/long_function_program.b");
	SET_DEST("benv/long_function_program2.b");
	goto(#replace_s);
	#change:
	change_flag = False; 
	SET_SOURCE("benv/long_function_program2.b");
	SET_DEST("benv/long_function_program.b");
	goto(#replace_s);

	#replace_e:
	del_file(change_flag);
};


void main(){

	int res; 

	res = init();

	[print(""), (0 == res), print("INIT ERROR\n")];

	replace();
	res = finish();

	[print(""), (0 == res), print("FINISH ERROR\n")];
};

main(); 



