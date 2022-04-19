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

int rindex(string s, string sub_s){
	int number;
	int s_len;
	int sub_len;
	int left_border;
	int res;
	
	number = index(s, sub_s);
	res = number;
	
	#rindex_s:
	[goto(#rindex_e), (-1 == number), print("")];
	s_len = len(s);
	sub_len = len(sub_s);
	left_border = (number + sub_len);
	s = s[left_border:s_len];
	number = index(s, sub_s);
	res = number;
	goto(#rindex_s);
	#rindex_e:
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
	string buf;
	string symbol;
	symbol="at";

	next_command(command);
	#get_imports_s:
	[goto(#get_imports_e), ("end" == command), print("")];
	number = index(command, "#import");
	[print(""), (0 == number), goto(#get_imports_e)];
	command_len = len(command);
	import_end = rindex(command, symbol);
	buf = str(import_end);
	print(buf);
	print("\n");

	import_name = command[7:command_len];
	res.push(import_name);
	next_command(command);	
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
