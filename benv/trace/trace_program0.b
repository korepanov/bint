$filestrings.b$;
$trace1$;
voidprintln(strings){print(s);
$trace2$;
print("\n");
$trace3$;
};
$fileio.b$;
$trace1$;
stackinput_data(){floata;
$trace2$;
floatb;
$trace3$;
floatc;
$trace4$;
stringbuf;
$trace5$;
println("Input a");
$trace6$;
input(buf);
$trace7$;
float$F0;
$F0=float(buf);
a=$F0;
UNDEFINE($F0);
$trace8$;
println("Input b");
$trace9$;
input(buf);
$trace10$;
float$F0;
$F0=float(buf);
b=$F0;
UNDEFINE($F0);
$trace11$;
println("Input c");
$trace12$;
input(buf);
$trace13$;
float$F0;
$F0=float(buf);
c=$F0;
UNDEFINE($F0);
$trace14$;
stackres;
$trace15$;
res.push(c);
$trace16$;
res.push(b);
$trace17$;
res.push(a);
$trace18$;
returnres;
$trace19$;
};
$trace20$;
voidoutput(floatx1,floatx2){stringbuf;
$trace21$;
println("Solution:");
$trace22$;
print("x1=");
$trace23$;
string$s0;
$s0=str(x1);
buf=$s0;
UNDEFINE($s0);
$trace24$;
println(buf);
$trace25$;
print("x2=");
$trace26$;
string$s0;
$s0=str(x2);
buf=$s0;
UNDEFINE($s0);
$trace27$;
println(buf);
$trace28$;
};
$filesolver.b$;
$trace1$;
floatdescr(floata,floatb,floatc){return((b^2.0)-((4.0*a)*c));
$trace2$;
};
$trace3$;
stacksolve(floata,floatb,floatc){floatx1;
$trace4$;
floatx2;
$trace5$;
stackres;
$trace6$;
floatd;
$trace7$;
d=descr(a,b,c);
$trace8$;
if(d<0){returnres;
$trace9$;
};
$trace10$;
x1=((((-1.0)*b)-(d^0.5))/(2.0*a));
$trace11$;
x2=((((-1.0)*b)+(d^0.5))/(2.0*a));
$trace12$;
res.push(x2);
$trace13$;
res.push(x1);
$trace14$;
returnres;
$trace15$;
};
$fileprogram.b$;
$trace1$;
voidmain(){println("Solving equation of the form ax^2 + bx + c = 0");
$trace2$;
stringbuf;
$trace3$;
do{stackres;
$trace4$;
res=input_data();
$trace5$;
floata;
$trace6$;
floatb;
$trace7$;
floatc;
$trace8$;
res.pop(a);
$trace9$;
res.pop(b);
$trace10$;
res.pop(c);
$trace11$;
res=solve(a,b,c);
$trace12$;
res.pop(buf);
$trace13$;
if("end"==buf){println("No solution");
$trace14$;
}else{floatx1;
$trace15$;
floatx2;
$trace16$;
float$F0;
$F0=float(buf);
x1=$F0;
UNDEFINE($F0);
$trace17$;
res.pop(x2);
$trace18$;
output(x1,x2);
$trace19$;
};
$trace20$;
println("again? y/n");
$trace21$;
input(buf);
$trace22$;
}while("y"==buf);
$trace23$;
};
$trace24$;
main();
