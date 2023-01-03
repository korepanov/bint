#import "stdlib/core.b"

string root_source;
string command;
bool e;

void init(){
	
	e = exists("benv/import_program.b");
	if (e){
		root_source = "benv/import_program.b";
	}else{
		root_source = "benv/trace_program.b";	
	};
	SET_SOURCE(root_source);
	SET_DEST("benv/slice_program.b");	
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
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
	int buf_pos;
	int epos;
	string find;
	string symbol;
	bool is_let;
	bool is_dig;

	find = "[";
	s = ops(command, find);
	
	string buf;
	s.pop(buf);
	while (NOT("end" == buf)){
		pos = int(buf);
		
		if(NOT(0 == pos)){
			buf_pos = (pos - 1);
			symbol = command[buf_pos];
			is_let = is_letter(symbol);
			is_dig = is_digit(symbol);
			
			if ((is_let)OR(is_dig)){	
				epos = slice_end(command, pos);
				el.push(epos);
				el.push(pos);
				res.push(el);
				el = null;	
			};
		};
		s.pop(buf);	
	};

	return res;
};

int colon_pos(string command){
	int res;
	string buf;
	stack temp;

	temp = ops(command, ":");
	temp.pop(buf);
	
	if (NOT("end" == buf)){
		res = int(buf);
		return res;
	};
	return -1;
};

int slice_name_start(string command, int slice_begin){
	stack el;
	int res;
	stack temp;
	
	command = command[0:slice_begin];
	temp = reg_find("[[:alpha:]]+[[:alnum:]]*$", command);
	temp.pop(el); 
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

void modify(){
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
			buf = ("int $sl_internal" + snumber);
			send_command(buf);
			buf = ("int $sl_left" + snumber);
			send_command(buf);
			buf = ("int $sl_right" + snumber);
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
};

void copy(string source, string dest){
	SET_SOURCE(source);
	SET_DEST(dest);
	string command;
	next_command(command); 
	while (NOT("end" == command)){
		send_command(command);
		next_command(command);	
	};

	UNSET_SOURCE();
	UNSET_DEST();
};

void main(){
	init();
	modify(); 
	finish();
	if ("benv/import_program.b" == root_source){
		copy("benv/slice_program.b", "benv/import_program.b");
	}else{
		copy("benv/slice_program.b", "benv/trace_program.b");
	};
	
	DEL_DEST("benv/slice_program.b");
	string t;
	string t2;
	string s;
	s = str(0);
	t = (("benv/trace/trace_program" + s) + ".b");
	e = exists(t);
	for (int number; number = 1; e; number = (number + 1)){
		t = (("benv/trace/trace_program" + s) + ".b");
		SET_SOURCE(t);
		t = (("benv/trace/slice_program" + s) + ".b");
		SET_DEST(t);
		modify();
		finish();
		t = (("benv/trace/slice_program" + s) + ".b");
		t2 = (("benv/trace/trace_program" + s) + ".b");
		copy(t, t2);
		t = (("benv/trace/slice_program" + s) + ".b"); 
		DEL_DEST(t);
		s = str(number);
		t = (("benv/trace/trace_program" + s) + ".b");
		e = exists(t);
	};
};

main();
