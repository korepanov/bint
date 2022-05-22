void println(string s){
	print(s);
	print("\n");
};

stack reverse(stack s){
	string buf;
	stack res; 
	
	s.pop(buf);
	#_reverse_s:
	[goto(#_reverse_e), ("end" == buf), print("")];
	res.push(buf);
	s.pop(buf);
	goto(#_reverse_s);
	#_reverse_e:
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
	#_indexes_s:
	[goto(#_indexes_e), (-1 == i), print("")];
	i = (i + (s_len_old - s_len));
	res.push(i);
	pointer = (pointer + sub_len);
	s = s[pointer:s_len];
	s_len = len(s);
	i = index(s, sub_s);
	pointer = i;
	goto(#_indexes_s);
	#_indexes_e:
	res = reverse(res);
	return res;
};

stack ops(string command, string op){
	stack quotes;
	stack these_quotes; 
	stack op_nums;
	string buf;
	stack res;
	int num1;
	int num2; 
	int op_num;
	bool was_quote;
	bool to_add;
	
	was_quote = False;
	to_add = True;

	op_nums = indexes(command, op);
	op_nums.pop(buf);
	#_op_nums_s:
	[goto(#_op_nums_e), ("end" == buf), print("")];
	op_num = int(buf);
	quotes = reg_find("\"(\\.|[^\"])*\"", command);
	
	#_quotes_s:
	quotes.pop(these_quotes);
	these_quotes.pop(buf);
	[goto(#_quotes_e), ("end" == buf), print("")];
	#_these_quotes_s:
	[goto(#_these_quotes_e), ("end" == buf), print("")];
	num1 = int(buf);
	these_quotes.pop(buf);
	num2 = int(buf);
	these_quotes.pop(buf);
	was_quote = True;
	goto(#_these_quotes_s);
	#_these_quotes_e:
	[print(""), ((op_num > num1)AND(op_num < num2)), goto(#_is_op_end)];
	to_add = False;
	goto(#_push_op_end);
	#_is_op_end:	
	goto(#_quotes_s);
	#_quotes_e:
	[goto(#_push_op_end), ((was_quote)AND(NOT(to_add))), print("")];
	res.push(op_num);
	was_quote = False;
	to_add = True;
	#_push_op_end:
	print("");
	op_nums.pop(buf);
	was_quote = False;
	to_add = True;
	goto(#_op_nums_s);
	#_op_nums_e:
	res = reverse(res);
	return res;
};
