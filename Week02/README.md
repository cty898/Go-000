学习笔记
作业批改记录：https://shimo.im/sheets/9AzpRIqEX4ssliEX/R5W0Z
参考答案：
dao:

 return errors.Wrapf(code.NotFound, fmt.Sprintf("sql: %s error: %v", sql, err))


biz:

if errors.Is(err, code.NotFound} {

}
