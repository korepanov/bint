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

void replace_print(string command){
	int quote_pos;
	int print_pos;

	quote_pos = index(command, "\"");
	print_pos = index(command, "print");
	
	[goto(#replace_print_e), (-1 == print_pos), print("")];
	[goto(#replace), (-1 == quote_pos), print("")];
	
	
	#replace:

	#replace_print_e:
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
