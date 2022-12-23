#import "stdlib/core.b"

string root_source;
string command;

void init(){
	root_source = "import_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/slice_program.b");	
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
	
	#braces_s:
	[print(""), (pos < command_len), goto(#braces_e)];
	spos = str(pos);
	[print(""), in_stack(obraces, spos), goto(#o_plus_end)];
	o_sum = (o_sum + 1);
	#o_plus_end:
	[print(""), in_stack(cbraces, spos), goto(#c_plus_end)];
	c_sum = (c_sum + 1);
	#c_plus_end:
	[goto(#braces_e), (o_sum == c_sum), print("")];
	pos = (pos + 1);
	goto(#braces_s);
	#braces_e:
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

int colon_pos(string command){
	int res;
	string buf;

	ops(command, ":").pop(buf);
	if (NOT("end" == buf)){
		res = int(buf);
		return res;
	};
	return -1;
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
	int colon;
	string buf;
	int buf_len;
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
		
		while (NOT("end" == buf)){
			
			s = slice_poses(command);
			
			s.pop(el);
			el.pop(buf);
			bpos = int(buf);
			el.pop(buf);
			epos = int(buf);

			gbpos = slice_name_start(command, bpos);
			bpos = (bpos + 1);
			snumber = str(number);
			buf = ("string $sl_internal" + snumber);
			send_command(buf);
			buf = ("string $sl_left" + snumber);
			send_command(buf);
			buf = ("string $sl_right" + snumber);
			send_command(buf);
			buf = command[bpos:epos];
			colon = colon_pos(buf);

			if (NOT(-1 == colon)){
				string buf2;
				buf2 = buf[0:colon];
				string buf3;
				colon = (colon + 1);
				buf_len = len(buf);
				buf3 = buf[colon:buf_len];
				buf = ((("$sl_left" + snumber) + "=") + buf2);
				send_command(buf);
				buf = ((("$sl_right" + snumber) + "=") + buf3);
				send_command(buf);
				sub_command = ((("$sl_left" + snumber) + ":") + ("$sl_right" + snumber));
			}else{
				buf = ((("$sl_internal" + snumber) + "=") + buf);
				send_command(buf);
				sub_command = ("$sl_internal" + snumber);
			};

			bpos = (bpos - 1);
			buf = ("string $sl" + snumber);
			send_command(buf);
			buf = command[gbpos:bpos];
			buf = (((((("$sl" + snumber) + "=") + buf) + "[") + sub_command) + "]");
			send_command(buf);
			sub_command = ("$sl" + snumber);
			epos = (epos + 1);
			command = modify_command(command, sub_command, gbpos, epos);
			number = (number + 1);
			s.pop(el);
			el.pop(buf);
		};
		send_command(command);

		for (int i; i = 0; i < number; i = (i + 1)){
			string b;
			snumber = str(i);
			b = (("UNDEFINE($sl_internal" + snumber) + ")");
			send_command(b);
			send_command((("UNDEFINE($sl" + snumber) + ")"));	
			send_command((("UNDEFINE($sl_left" + snumber) + ")"));
			send_command((("UNDEFINE($sl_right" + snumber) + ")"));		
		};

		next_command(command);
	};
	finish();
};

main();
