// Referencia em C# para implementação interna do sistema de pulse:
public delegate void MyEventHandler();

public class MyClass
{
    public event MyEventHandler BeforeAction;
    public event MyEventHandler AfterAction;

    public void DoAction()
    {
        if (BeforeAction != null) BeforeAction();
        // Realiza a ação aqui
        if (AfterAction != null) AfterAction();
    }
}

public class Program
{
    public static void Main()
    {
        MyClass myClass = new MyClass();
        myClass.BeforeAction += () => Console.WriteLine("Before: Action is about to be executed.");
        myClass.AfterAction += () => Console.WriteLine("After: Action has been executed.");

        myClass.DoAction();
    }
}
