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
	stack indexes;
	stack these_indexes; 
	string buf;
	stack res;

	indexes = reg_find("\"(.*?)\"", command);
	print(command);
	print("\n");
	#indexes_s:
	indexes.pop(these_indexes);
	these_indexes.pop(buf);
	[goto(#indexes_e), ("end" == buf), print("")];
	#these_indexes_s:
	[goto(#these_indexes_e), ("end" == buf), print("")];
	print(buf);
	print(" ");
	these_indexes.pop(buf);
	goto(#these_indexes_s);
	#these_indexes_e:
	print("\n");
	goto(#indexes_s);
	#indexes_e:

	return res;
};

void replace_print(string command){
	stack res;
	res = ops(command, "print");
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
