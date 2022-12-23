#import "stdlib/core.b"

string root_source;
string command;

void init(){
	root_source = "import_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/index_program.b");	
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	//DEL_DEST(root_source); 
};

int slice_end(string command, int func_begin){
	stack obraces;
	stack cbraces;
	string obrace;
	string cbrace;
	string symbol;
	int o_sum;
	int c_sum;
	int pos;
	string spos;
	int command_len;
	
	command_len = len(command);
	obrace = "[";
	cbrace = "]";
	o_sum = 1;
	c_sum = 0;
	pos = (func_begin + 1);

	obraces = ops(command, obrace);
	cbraces = ops(command, cbrace);
	
	#_braces_s:
	[print(""), (pos < command_len), goto(#_braces_e)];
	spos = str(pos);
	[print(""), in_stack(obraces, spos), goto(#_o_plus_end)];
	o_sum = (o_sum + 1);
	#_o_plus_end:
	[print(""), in_stack(cbraces, spos), goto(#_c_plus_end)];
	c_sum = (c_sum + 1);
	#_c_plus_end:
	[goto(#_braces_e), (o_sum == c_sum), print("")];
	pos = (pos + 1);
	goto(#_braces_s);
	#_braces_e:
	return pos;
};

stack slice_poses(string command){
	stack s;
	stack el;
	stack res;
	stack null; 
	int pos;
	int epos;
	string find;

	find = "[";
	s = ops(command, find);
	string buf;
	s.pop(buf);
	
	while (NOT("end" == buf)){
		pos = int(buf);
		if(NOT(0 == pos)){
			epos = slice_end(command, pos);
			el.push(epos);
			el.push(pos);
			res.push(el);
			el = null;	
		};
		s.pop(buf);	
	};

	return res;
};

int slice_name_start(string command, int slice_begin){
	stack el;
	int res;

	command = command[0:slice_begin];
	reg_find("[[:alpha:]]+[[:alnum:]]*$", command).pop(el);
	el.pop(res);
	return res; 
};

string modify_command(string command, string sub_command, int bpos, int epos){
	string new_command;
	string buf;
	int command_len;
 
	new_command = command[0:bpos];
	new_command = (new_command + sub_command);
	command_len = len(command);
	buf = command[epos:command_len];
	new_command = (new_command + buf);
	
	return new_command;
};

void main(){
	stack s;
	stack el;
	int bpos;
	int epos;
	int gbpos;
	int command_len;
	string buf;
	int number;
	string snumber;
	string sub_command; 

	init();
	next_command(command);
	
	while (NOT("end" == command)){
		number = 0;
		s = slice_poses(command);
		
		s.pop(el);
		el.pop(buf);
		print((command  + "\n"));
		while (NOT("end" == buf)){
			
			s = slice_poses(command);
			
			s.pop(el);
			el.pop(buf);
			bpos = int(buf);
			el.pop(buf);
			epos = int(buf);
			epos = (epos + 1);

			gbpos = slice_name_start(command, bpos);
			print((command + "\n")); 
			command_len = len(command);
			command = command[gbpos:command_len];
			print((command + "\n"));
		
			snumber = str(number);
			buf = ("int $sl" + snumber);
			send_command(buf);
			buf = command[bpos:epos];
			buf = ((("$sl" + snumber) + "=") + buf);
			send_command(buf);
			sub_command = ("$sl" + snumber);
			command = modify_command(command, sub_command, bpos, epos);
			number = (number + 1);
			s.pop(el);
			el.pop(buf);
		};
		send_command(command);

		for (int i; i = 0; i < number; i = (i + 1)){
			string b;
			snumber = str(i);
			b = (("UNDEFINE($sl" + snumber) + ")");
			send_command(b);		
		};

		next_command(command);
	};
	finish();
};

main();
