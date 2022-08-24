#import "stdlib/core.b"
bool first_file; 
int br_closed;
int br_opened;

void init(){
	string root_source;
	first_file = True;
	br_closed = 0;
	br_opened = 0;
	root_source = "benv/long_function_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/recurs_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_var_def(string command){
	stack s;
	string buf;

	s = reg_find("^(?:(int|float|bool|stack|string)[^\(]*)", command);
	s.pop(buf);
	
	return (NOT("end" == buf)); 
};

bool is_func_definition(string command){
	stack s;
	string buf;

	string r1;
	string r2;
	string r3;
	string r;

	r1 = "(?:(int|float|bool|string|stack|void)[[:alnum:]|_]";
	r2 = "*?\((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})";
	r3 = "(int|float|bool|string|stack)[[:alnum:]|_]+\){";
	
	r = ((r1 + r2) + r3);

	s = reg_find(r, command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return True;
	};	

	s = reg_find("(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\(\){)", command);
	
	s.pop(buf);
	return (NOT("end" == buf));
};

string T(string command){
	stack s;
	string buf;

	s = reg_find("(?:(^int))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "int";	
	};

	s = reg_find("(?:(^float))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "float";	
	};

	s = reg_find("(?:(^bool))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "bool";	
	};

	s = reg_find("(?:(^stack))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "stack";	
	};

	s = reg_find("(?:(^string))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "string";	
	};

	s = reg_find("(?:(^void))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "void";	
	};
	
	print(command);
	print("\n");
	print("T: ERROR\n");

	return "ERROR";	
};

string func_name(string command){
	string t;
	int pos;
	int epos;
	int command_len;

	t = T(command);
	if ("ERROR" == command){
		return "ERROR";	
	};
	
	pos = len(t);
	command_len = len(command);
	command = command[pos:command_len];
	
	epos = index(command, "(");

	if (-1 == epos){
		return "ERROR";	
	};

	command = command[0:epos];

	return command;

};

stack args_to_accept(string command){
	int pos1;
	int pos2;
	string s_to_parse;
	string t;
	stack res;
	stack null;
	stack el;
	int L;
	string arg_name;
	string arg_type;
	
	pos1 = index(command, "(");
	pos1 = (pos1 + 1);
	pos2 = index(command, ")");

	if (pos1 == pos2){
		return res; 	
	};

	if ((-1 == pos1) OR (-1 == pos2)){
		print(command);
		print("\n");
		print("recurs: args_to_accept: ERROR\n");	
	};

	s_to_parse = command[pos1:pos2];
	
	#parse_s:	
	el = null;
	t = T(s_to_parse);
	pos1 = len(t);
	L = len(s_to_parse);

	s_to_parse = s_to_parse[pos1:L];
	
	pos1 = index(s_to_parse, ",");
	if (-1 == pos1){
		pos1 = len(s_to_parse);	
		arg_name = s_to_parse[0:pos1];
		el.push(arg_name);
		el.push(t);
		res.push(el);

		stack res_el;
		stack new_res;
		string buf;

		#rev:
		res.pop(res_el);
		res_el.pop(buf);

		if ("end" == buf){
			return new_res;		
		};

		res_el.push(buf);
		new_res.push(res_el);		

		goto(#rev); 
	};
	
	arg_name = s_to_parse[0:pos1];
	pos1 = (pos1 + 1);
	L = len(s_to_parse);
	s_to_parse = s_to_parse[pos1:L];
	el.push(arg_name);
	el.push(t);
	res.push(el);

	goto(#parse_s);
};

stack args_to_pass(string fcall){
	int pos1;
	int pos2;
	int command_len;
	string buf;
	stack res; 
	string command;

	command = fcall;

	pos1 = index(command, "(");
	pos2 = func_end(command, 0);
	command = command[pos1:pos2];
	pos1 = 1;
	command_len = len(command);
	
	if (pos1 == (command_len - 1)){
		return res;	
	};

	#next_arg:
	pos2 = index(command, ",");

	if (-1 == pos2){
		command_len = len(command);
		pos2 = (command_len - 1);
		buf = command[pos1:pos2];
		res.push(buf);
		return res;
	};

	buf = command[pos1:pos2];
	res.push(buf);
	pos2 = (pos2 + 1);
	command_len = len(command);
	command = command[pos2:command_len];
	pos1 = 0;
	goto(#next_arg);
};

int func_call(string fname, string command){
	string reg;
	stack s; 
	stack this_s;
	string buf;
	int pos;
	
	reg = (("(?:" + fname) + "\(.*\))");
	s = reg_find(reg, command);
	s.pop(this_s);
	this_s.pop(buf);
	
	if ("end" == buf){
		pos = -1;
	}else{
		pos = int(buf);
	};
	
	
	return pos;
};

void switch_files(){
	finish();
	[print(""), (first_file), goto(#first_end)];
	SET_SOURCE("benv/recurs_program.b");
	SET_DEST("benv/recurs_program2.b");
	first_file = False;
	goto(#switch_files_e);
	#first_end:
	SET_SOURCE("benv/recurs_program2.b");
	SET_DEST("benv/recurs_program.b");
	first_file = True; 
	#switch_files_e:
	print("");
};

void clear_files(){
	finish();
	[goto(#first), (first_file), print("")];
	SET_SOURCE("benv/recurs_program2.b");
	SET_DEST("benv/recurs_program.b");
	switch_files();
	switch_command();
	#clear_files_s:
	[goto(#clear_files_e), ("end" == command), print("")];
	send_command(command);
	switch_command();
	goto(#clear_files_s);

	#first:
	print("");
	#clear_files_e:
	DEL_DEST("benv/recurs_program2.b");
};

void check_br(string command){
	string symbol;
	stack s;
	string buf;

	symbol = "{";
	s = ops(command, symbol);
	s.pop(buf);
	
	if (NOT("end" == buf)){
		br_opened = (br_opened + 1);
	};

	symbol = "}";
	s = ops(command, symbol);
	s.pop(buf);
	
	if (NOT("end" == buf)){
		br_closed = (br_closed + 1);
	};
};

void reset_br(){
	br_opened = 0;
	br_closed = 0;
};

void main(){
	init();
	string name;
	stack accepted_args;
	stack passed_args;
	stack raccepted_args;
	stack args_to_undefine;
	stack null; 
	stack buf_stack;
	string t;
	stack arg;
	string arg_name;
	string arg_type;
	int counter;
	string buf;
	int old_COMMAND_COUNTER;
	int buf_COMMAND_COUNTER;
	int pos;
	int pos1;
	int pos2;
	int call_counter;
	int old_call_counter;
	string new_command;
	bool was_recurs;
	bool first_recurs_call;
	int call_num;
	string scall_num;
	string res;
	int command_len;
	string command_buf;
	string original_command;

	call_counter = 0;
	was_recurs = False;
	first_recurs_call = True;
	send_command("string $ret");
	switch_command();
	
	#main_s:
	print("");
	if (NOT("end" == command)){
		if (is_func_definition(command)){
			t = T(command);
			if ("ERROR" == t){
				print("recurs ERROR\n");			
			};
			name = func_name(command);
			if ("ERROR" == name){
				print("recurs ERROR\n");			
			};
			call_num = 0;
			raccepted_args = args_to_accept(command);

			counter = block_end();
			old_COMMAND_COUNTER = COMMAND_COUNTER;

			#is_recurs:
			print("");
			if (COMMAND_COUNTER < counter){
				switch_command();
				pos = func_call(name, command);
				if (NOT(-1 == pos)){
					was_recurs = True;
					old_call_counter = call_counter;
					call_counter = COMMAND_COUNTER;
					SET_COMMAND_COUNTER(old_call_counter);
					
					switch_command();
					#block_begin:
					print("");
					if (COMMAND_COUNTER < old_COMMAND_COUNTER){
						send_command(command);
						switch_command();
						goto(#block_begin);
					};
					
					if (first_recurs_call){
						first_recurs_call = False;
						new_command = ("stack" + ("$" + (name + "_stack")));
						send_command(new_command);
						new_command = (("goto(#" + name) + "_end)");
						send_command(new_command);
						new_command = (("#" + name) + ":");
						new_command = (new_command + "print(\"\")");
						send_command(new_command); 

						accepted_args = raccepted_args;
						
						accepted_args.pop(arg);			
						arg.pop(arg_type);
						arg.pop(arg_name);

						#append_args:
						if (NOT("end" == arg_name)){
							args_to_undefine.push(arg_name);
							accepted_args.pop(arg);			
							arg.pop(arg_type);
							arg.pop(arg_name);
							goto(#append_args);		
						};

						switch_command();
						
					};
					
					accepted_args = raccepted_args;
					accepted_args.pop(arg);
					arg.pop(arg_type);
					arg.pop(arg_name);
					#ax:
					if (NOT("end" == arg_name)){
						new_command = (arg_type + arg_name);
						send_command(new_command);
						new_command = (("pop(" + arg_name) + ")");
						send_command(new_command);
						accepted_args.pop(arg);
						arg.pop(arg_type);
						arg.pop(arg_name);
						goto(#ax); 					
					};
					bool is_ret;
					is_ret = False;
					#before_recurs:
					print("");
					if (COMMAND_COUNTER < call_counter){
						command_len = len(command);
						if (command_len > 6){
							command_buf = command[0:6];
							if ("return" == command_buf){
								string ret_arg;
								ret_arg = command[6:command_len];
								new_command = (("push(" + ret_arg) + ")");
								send_command(new_command);

								buf_stack = args_to_undefine;
								args_to_undefine.pop(arg_name);
								#un:
								if (NOT("end" == arg_name)){
									new_command = (("UNDEFINE(" + arg_name) + ")");
									send_command(new_command);
									args_to_undefine.pop(arg_name);
									goto(#un);
								};
								args_to_undefine = buf_stack;
		
								new_command = (("$" + name) + "_stack.pop($ret)");
								send_command(new_command);
								new_command = "goto($ret)";
								send_command(new_command); 	
								is_ret = True;						
							};	
												
						};

						check_br(command);
						
						if ((is_var_def(command))AND(br_closed == br_opened)){
							arg_type = T(command);
							int type_len; 
							type_len = len(arg_type);
							command_len = len(command);
							arg_name = command[type_len:command_len];
							args_to_undefine.push(arg_name);
						};

						if (NOT(is_ret)){
							send_command(command);
						};
						is_ret = False;
						switch_command();
						goto(#before_recurs);
					};
					scall_num = str(call_num);
					res = ((("#" + name) + "_res") + scall_num);
					new_command = ((("$" + (name + "_stack.push(")) + res) + ")");
					send_command(new_command);
					
					call_num = (call_num + 1);
					pos1 = func_call(name, command);
					pos2 = func_end(command, pos1);
	
					new_command = command[pos1:pos2];
					passed_args = args_to_pass(new_command);
					
					passed_args.pop(arg_name);

					#internal_passed:
					print("");
					if (NOT("end" == arg_name)){
						res = (("push(" + arg_name) + ")");
						send_command(res);
						passed_args.pop(arg_name);
						goto(#internal_passed);						
					};
					
					new_command = (("goto(#" + name) + ")");
					send_command(new_command);
					res = ((("#" + name) + "_res") + scall_num);
					res = (res + ":print(\"\")");
					send_command(res);

					if (NOT("void"==t)){
						new_command = ((((t + "$") + name) + "_recurs_res") + scall_num); 
						send_command(new_command);
						new_command = (("pop(" + ((("$" + name) + "_recurs_res") + scall_num)) + ")");
						send_command(new_command);
						new_command = command[0:pos1];
						new_command = ((((new_command + "$") + name) + "_recurs_res") + scall_num);
						command_len = len(command);
						command = command[pos2:command_len];
						new_command = (new_command + command);
					

						command_len = len(new_command);
						
						if (command_len > 6){
							command = new_command[0:6];
							
							if ("return" == command){
								string ret_arg;
								ret_arg = new_command[6:command_len];
								command = (("push(" + ret_arg) + ")");
								send_command(command);

								buf_stack = args_to_undefine;
								args_to_undefine.pop(arg_name);
								#un2:
								if (NOT("end" == arg_name)){
									new_command = (("UNDEFINE(" + arg_name) + ")");
									send_command(new_command);
									args_to_undefine.pop(arg_name);
									goto(#un2);
								};
								args_to_undefine = buf_stack;
								
								command = (("$" + name) + "_stack.pop($ret)");
								send_command(command);
								command = "goto($ret)";
								send_command(command); 
							}else{
								send_command(new_command);
							};				
						}else{
							send_command(new_command);
						};
					};
					
				};
				goto(#is_recurs);
			};
			
			if (was_recurs){
				was_recurs = False; 
				first_recurs_call = True;
				
				SET_COMMAND_COUNTER(call_counter);
				
				switch_command();
				
				bool is_ret;
				is_ret = False;
				#rest:	
				print("");
				if (COMMAND_COUNTER < counter){
					command_len = len(command);
					if (command_len > 6){
						command_buf = command[0:6];
						if ("return" == command_buf){
							string ret_arg;
							ret_arg = command[6:command_len];
							new_command = (("push(" + ret_arg) + ")");
							send_command(new_command);
							
							buf_stack = args_to_undefine;
							args_to_undefine.pop(arg_name);
							#un4:
							if (NOT("end" == arg_name)){
								new_command = (("UNDEFINE(" + arg_name) + ")");
								send_command(new_command);
								args_to_undefine.pop(arg_name);
								goto(#un4);
							};
							args_to_undefine = buf_stack;

							new_command = (("$" + name) + "_stack.pop($ret)");
							send_command(new_command);
							new_command = "goto($ret)";
							send_command(new_command);
							is_ret = True; 							
						};							
					};
					
					check_br(command);
					
					if ((is_var_def(command))AND(br_closed == br_opened)){
						arg_type = T(command);
						int type_len; 
						type_len = len(arg_type);
						command_len = len(command);
						arg_name = command[type_len:command_len];
						args_to_undefine.push(arg_name);
					};

					if (NOT(is_ret)){
						send_command(command);
					};
					is_ret = False;
					switch_command();
					goto(#rest);
					
				}; 
				if ("void" == t){
					args_to_undefine.pop(arg_name);
					#un3:
					if (NOT("end" == arg_name)){
						new_command = (("UNDEFINE(" + arg_name) + ")");
						send_command(new_command);
						args_to_undefine.pop(arg_name);
						goto(#un3);
					};		
					
					new_command = (("$" + name) + "_stack.pop($ret)");
					send_command(new_command);
					new_command = "goto($ret)";
					send_command(new_command);		
				};
				new_command = ("#" + (name + "_end:"));
				new_command = (new_command + "print(\"\")");
				send_command(new_command);				

				switch_command();
				
				#external_call:
				print("");
				if(NOT("end" == command)){
					pos = func_call(name, command);
					if (NOT(-1 == pos)){
						pos1 = pos;
						pos2 = func_end(command, pos1);
						original_command = command;
						command = command[pos1:pos2];

						passed_args = args_to_pass(command);
						
						passed_args.pop(arg_name);

						#passed:
						print("");
						if (NOT("end" == arg_name)){
							res = (("push(" + arg_name) + ")");
							send_command(res);
							passed_args.pop(arg_name);
							goto(#passed);						
						};

						scall_num = str(call_num);
						res = ((((((("$" + name) + "_stack.push(") + "#") + name) + "_res") + scall_num) + ")");
						send_command(res);
						res = (("goto(#" + name) + ")");
						send_command(res); 
						res = ((("#" + name) + "_res") + scall_num);
						res = (res + ":print(\"\")");
						send_command(res);

						if (NOT("void" == t)){
							res = (t + ((("$" + name) + "_recurs_res") + scall_num));
							send_command(res);
							res = (((("pop($" + name) + "_recurs_res") + scall_num) + ")");
							send_command(res);
							command = original_command[0:pos1];
							command = (command + ((("$" + name) + "_recurs_res") + scall_num));
							command_len = len(original_command);
							original_command = original_command[pos2:command_len];
							command = (command + original_command);
							send_command(command);
						};

						call_num = (call_num + 1);
					};

					if (-1 == pos){
						send_command(command);
					};

					switch_command();
					goto(#external_call);
				};
				switch_files();
				call_counter = 0;
				SET_COMMAND_COUNTER(0);
			};
			
		};
		reset_br();
		args_to_undefine = null; 
		switch_command();
		goto(#main_s);
	};
	
	SET_COMMAND_COUNTER(0);
	switch_command();
	#final_rest:
	print("");
	if (NOT("end" == command)){
		send_command(command);
		switch_command();
		goto(#final_rest);
	};

	clear_files();
	
};

main();
