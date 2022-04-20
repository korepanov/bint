string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/import_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

int find_import_end(string command){
	int command_len;
	int res;

	command_len = len(command);	
	command = command[8:command_len];
	res = index(command, "\"");
	res = (res + 8);	
 
	return res;
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

stack get_imports(){
	stack res; 
	int number;
	string import_name;
	int command_len;
	int import_end;
	
	next_command(command);
	#get_imports_s:
	[goto(#get_imports_e), ("end" == command), print("")];
	number = index(command, "#import");
	[print(""), (0 == number), goto(#get_imports_e)];
	command_len = len(command);
	import_end = find_import_end(command);
	import_name = command[8:import_end];
	res.push(import_name);
	import_end = (import_end + 1);
	command = command[import_end:command_len];
	goto(#get_imports_s);
	#get_imports_e:
	RESET_SOURCE();
	return reverse(res);	
};

void file_union(){
	imports = get_imports();
	imports.pop(import);
	#file_union_s:
	[goto(#file_union_e), ("end" == import), print("")];
	print(import);
	print("\n");
	imports.pop(import);	
	goto(#file_union_s);
	#file_union_e:
	print("");
};

void main(){
	string command;
	stack imports;
	string import;
	
	init();
	file_union();	
	finish();
};

main();
