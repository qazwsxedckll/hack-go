@17
D=A
@SP
A=M
M=D
@SP
M=M+1
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_0
D;JEQ
@SP
A=M-1
M=0
@END_0
0;JMP
(JUMP_0)
@SP
A=M-1
M=-1
(END_0)
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
@16
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_1
D;JEQ
@SP
A=M-1
M=0
@END_1
0;JMP
(JUMP_1)
@SP
A=M-1
M=-1
(END_1)
@16
D=A
@SP
A=M
M=D
@SP
M=M+1
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_2
D;JEQ
@SP
A=M-1
M=0
@END_2
0;JMP
(JUMP_2)
@SP
A=M-1
M=-1
(END_2)
@892
D=A
@SP
A=M
M=D
@SP
M=M+1
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_3
D;JLT
@SP
A=M-1
M=0
@END_3
0;JMP
(JUMP_3)
@SP
A=M-1
M=-1
(END_3)
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
@892
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_4
D;JLT
@SP
A=M-1
M=0
@END_4
0;JMP
(JUMP_4)
@SP
A=M-1
M=-1
(END_4)
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_5
D;JLT
@SP
A=M-1
M=0
@END_5
0;JMP
(JUMP_5)
@SP
A=M-1
M=-1
(END_5)
@32767
D=A
@SP
A=M
M=D
@SP
M=M+1
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_6
D;JGT
@SP
A=M-1
M=0
@END_6
0;JMP
(JUMP_6)
@SP
A=M-1
M=-1
(END_6)
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
@32767
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_7
D;JGT
@SP
A=M-1
M=0
@END_7
0;JMP
(JUMP_7)
@SP
A=M-1
M=-1
(END_7)
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
D=M-D
@JUMP_8
D;JGT
@SP
A=M-1
M=0
@END_8
0;JMP
(JUMP_8)
@SP
A=M-1
M=-1
(END_8)
@57
D=A
@SP
A=M
M=D
@SP
M=M+1
@31
D=A
@SP
A=M
M=D
@SP
M=M+1
@53
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
M=M+D
@112
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
M=M-D
@SP
AM=M-1
M=-M
@SP
AM=M-1
D=M
A=A-1
M=M&D
@82
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
AM=M-1
D=M
A=A-1
M=M|D
@SP
AM=M-1
M=!M
