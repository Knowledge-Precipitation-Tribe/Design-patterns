package 装饰者模式;

/**
 * Created by ${super} on ${2017/5/27}.
 */
public class TShirts extends Finery {
    @Override
    public void show() {
        System.out.println("大T恤");
        super.show();
    }
}
