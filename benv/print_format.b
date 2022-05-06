string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/print_format_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

stack ops(string command, string op){
	stack res;
	int i;
	int op_len;
	int command_len;
	string lexeme;

	command_len = len(command);
	op_len = len(op);
	i = 0;

	#ops_s:
	[print(""), (i < command_len), goto(#ops_e)];
	[goto(#push_op_s), (op == command[i:op_len]), print("")];
	[goto(#qoute_s), ("\"" == command[i]), ];
	#quote_s:
	
	#quote_e:
	#push_op_s:
	i = (i + command_len);
	goto(#ops_s);
	#push_op_e:
	#ops_e: 

	return res;
};

void replace_print(string command){
	
	print("");
};

void main(){
	string command;
	init();
	
	next_command(command);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	replace_print(command);
	send_command(command);
	next_command(command);
	goto(#main_s);
	#main_e:
	finish();
};

main();
